package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

const testPairs = 10

// newTestDB points the global db at a fresh SQLite file holding testPairs pairs.
func newTestDB(t *testing.T) {
	t.Helper()
	db.init(filepath.Join(t.TempDir(), "test.db"))
	t.Cleanup(db.Close)

	_, err := db.sqldb.Exec(`
		CREATE TABLE pairs (
			id INTEGER PRIMARY KEY,
			left TEXT NOT NULL,
			right TEXT NOT NULL,
			lcount INTEGER DEFAULT 0 NOT NULL,
			rcount INTEGER DEFAULT 0 NOT NULL
		)`)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= testPairs; i++ {
		if _, err := db.sqldb.Exec("INSERT INTO pairs (left, right) VALUES (?, ?)", fmt.Sprint("left", i), fmt.Sprint("right", i)); err != nil {
			t.Fatal(err)
		}
	}
	NUMBER_OF_PAIRS = db.getNumberOfPairs()
}

func get(path, user string) *httptest.ResponseRecorder {
	handler := map[string]http.HandlerFunc{
		"/random-pair":    getRandomPairHandler,
		"/n-random-pairs": handleGetNumberOfPairsN,
	}[path]
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if user != "" {
		req.AddCookie(&http.Cookie{Name: "user_id", Value: user})
	}
	rec := httptest.NewRecorder()
	handler(rec, req)
	return rec
}

func vote(user string, id int, side string) *httptest.ResponseRecorder {
	body := fmt.Sprintf(`{"id":%d,"leftright":%q}`, id, side)
	req := httptest.NewRequest(http.MethodPost, "/store-answer", strings.NewReader(body))
	if user != "" {
		req.AddCookie(&http.Cookie{Name: "user_id", Value: user})
	}
	rec := httptest.NewRecorder()
	storeAnswer(rec, req)
	return rec
}

// fetchIDs makes one request to path as user and returns the served pair IDs,
// and whether the server reported that every pair has been seen.
func fetchIDs(t *testing.T, path, user string) (ids []int, allSeen bool) {
	t.Helper()
	rec := get(path, user)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET %s: status %d: %s", path, rec.Code, rec.Body)
	}
	if path == "/random-pair" {
		var res Response
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		return []int{res.Pair.Id}, res.AllPairsSeen
	}
	var res MultiPairResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}
	for _, p := range res.Pairs {
		ids = append(ids, p.Id)
	}
	return ids, res.AllPairsSeen
}

// servedUntilDone fetches from path as user until the server reports that
// every pair has been seen, and returns the IDs served on the way.
func servedUntilDone(t *testing.T, path, user string) []int {
	t.Helper()
	var served []int
	for range testPairs + 1 {
		ids, allSeen := fetchIDs(t, path, user)
		if allSeen {
			return served
		}
		served = append(served, ids...)
	}
	t.Fatalf("%s never reported allPairsSeen for %s; served %v", path, user, served)
	return nil
}

func assertEveryPairOnce(t *testing.T, ids []int) {
	t.Helper()
	got := slices.Sorted(slices.Values(ids))
	want := make([]int, testPairs)
	for i := range want {
		want[i] = i + 1
	}
	if !slices.Equal(got, want) {
		t.Errorf("served pairs %v, want each of 1..%d exactly once", got, testPairs)
	}
}

func votesFor(t *testing.T, id int) int {
	t.Helper()
	var n int
	if err := db.sqldb.QueryRow("SELECT lcount + rcount FROM pairs WHERE id = ?", id).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func totalVotes(t *testing.T) int {
	t.Helper()
	var n int
	if err := db.sqldb.QueryRow("SELECT SUM(lcount + rcount) FROM pairs").Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func TestConcurrentUsersAllVotesCounted(t *testing.T) {
	newTestDB(t)

	var wg sync.WaitGroup
	var votes atomic.Int64
	for i := range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			user := fmt.Sprint("concurrent-", i)
			var res MultiPairResponse
			if err := json.NewDecoder(get("/n-random-pairs", user).Body).Decode(&res); err != nil {
				t.Error(err)
				return
			}
			for _, p := range res.Pairs {
				if rec := vote(user, p.Id, "left"); rec.Code != http.StatusOK {
					t.Errorf("vote by %s on pair %d: status %d: %s", user, p.Id, rec.Code, rec.Body)
					continue
				}
				votes.Add(1)
			}
		}()
	}
	wg.Wait()

	if got, want := totalVotes(t), int(votes.Load()); got != want {
		t.Errorf("stored %d votes, want %d", got, want)
	}
}

func TestFinishingDoesNotResetOtherUsers(t *testing.T) {
	newTestDB(t)

	seenByB, _ := fetchIDs(t, "/n-random-pairs", "user-b")
	servedUntilDone(t, "/n-random-pairs", "user-a")
	seenByB = append(seenByB, servedUntilDone(t, "/n-random-pairs", "user-b")...)

	assertEveryPairOnce(t, seenByB)
}

func TestSecondCycleServesEveryPair(t *testing.T) {
	newTestDB(t)

	for _, path := range []string{"/random-pair", "/n-random-pairs"} {
		t.Run(path, func(t *testing.T) {
			user := "cycler" + path
			assertEveryPairOnce(t, servedUntilDone(t, path, user))
			assertEveryPairOnce(t, servedUntilDone(t, path, user))
		})
	}
}

func TestVoteWithoutCookieIsRejected(t *testing.T) {
	newTestDB(t)

	if rec := vote("", 1, "left"); rec.Code != http.StatusBadRequest {
		t.Errorf("status %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if n := votesFor(t, 1); n != 0 {
		t.Errorf("pair 1 has %d votes, want 0", n)
	}
}

func TestVoteFromUnknownUserIsNotCounted(t *testing.T) {
	newTestDB(t)

	if rec := vote("stranger", 1, "left"); rec.Code != http.StatusBadRequest {
		t.Errorf("status %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if n := votesFor(t, 1); n != 0 {
		t.Errorf("pair 1 has %d votes, want 0", n)
	}
}

func TestVoteAfterFinishingIsCounted(t *testing.T) {
	newTestDB(t)

	servedUntilDone(t, "/n-random-pairs", "finisher")

	if rec := vote("finisher", 1, "right"); rec.Code != http.StatusOK {
		t.Errorf("status %d, want %d: %s", rec.Code, http.StatusOK, rec.Body)
	}
	if n := votesFor(t, 1); n != 1 {
		t.Errorf("pair 1 has %d votes, want 1", n)
	}
}

func TestFetchReportsDatabaseErrors(t *testing.T) {
	newTestDB(t)
	if _, err := db.sqldb.Exec("DROP TABLE pairs"); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{"/random-pair", "/n-random-pairs"} {
		if rec := get(path, "user"+path); rec.Code != http.StatusInternalServerError {
			t.Errorf("GET %s: status %d, want %d", path, rec.Code, http.StatusInternalServerError)
		}
	}
}
