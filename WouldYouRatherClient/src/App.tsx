import { useEffect, useState } from "react"
import { NavBar } from "./components/NavBar"
import { DisplayMessage } from "./components/DisplayMessage"
import { CardWrapper } from "./components/CardWrapper"
import Pair from "./types/Pair"

const random_pair_url = import.meta.env.VITE_RANDOM_PAIR_URL
const n_random_pairs_url = import.meta.env.VITE_N_RANDOM_PAIRS_URL


function App() {
  const [pair,setPair] = useState<Pair>({id:-1, left:"",right:""})
  const [unseenPairs, setUnseenPairs] = useState<Pair[]>([])

  enum State {
    START = 1,
    PLAY = 2,
    ANSWERED =3,
    END = 4,
  }

  const [gameState, setGameState] = useState(State.START)

  useEffect(() => {
    fetchNPairsTESTING()
    const savedTheme = localStorage.getItem("theme")
    if(savedTheme) {
      document.documentElement.setAttribute('data-theme', savedTheme)
    }
  },[])

  function handleChangeTheme(newTheme:string) {
    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
  }

  //fetch new pair from server
  async function fetchPair(){
    setPair({id:-1, left:"",right:""})
    setGameState(State.PLAY)
    const res = await fetch(random_pair_url, {method:"GET", credentials:"include",headers: {"Content-Type":"application/json"}})
    if(res.ok) {
      const data = await res.json()
      if(data.allPairsSeen){
        setGameState(State.END)
      }
      setPair({id:data.pair.id, left:data.pair.left, right:data.pair.right})
    }
    else{
      console.log("API call failed", res.status)
    }
  }
  ///TESTING HERE
  async function fetchNPairsTESTING(){
    setPair({id:-1, left:"",right:""})
    const res = await fetch(n_random_pairs_url, {method:"GET", credentials:"include",headers: {"Content-Type":"application/json"}})
    if(res.ok) {
      const data = await res.json()
      console.log(data)
      const pairs:Pair[] = data.pair.map((item: any) =>({id:item.id, left:item.left, right:item.right}))
      setUnseenPairs((unseenPairs) => [...unseenPairs, ...pairs])
    }
    else{
      console.log("API call failed", res.status)
    }
  }

  function handlePlayAgain(){
    setGameState(State.START)
    setPair({id:-1, left:"",right:""})
  }

  function handleAnswer(e: React.MouseEvent<HTMLDivElement, MouseEvent>){
    e.preventDefault()
    setGameState(State.ANSWERED)
  }

  function handleOnPlay() {
    setGameState(State.PLAY)
    const tempPairs = [...unseenPairs]
    const pair = tempPairs.pop()
    if(pair !== undefined){
      setPair(pair)
    }
    if(unseenPairs.length < 2){
      fetchNPairsTESTING()
    }
    else{
      setUnseenPairs(tempPairs)
    }
  }

  const mainContent = () => {
    switch(gameState){ 
      case State.START:
        return (<DisplayMessage headingText="Welcome to Would you rather, press start to play"></DisplayMessage>)
      case State.PLAY:
        return (<CardWrapper handleAnswer={handleAnswer} pair={pair}></CardWrapper>)
      case State.ANSWERED:
        return (<CardWrapper pair={pair}></CardWrapper>)
      case State.END:
        return (<DisplayMessage handleClick={handlePlayAgain} buttonText="Play again!" headingText="You have finished all questions! come back another time or"></DisplayMessage>)
      }
  } 

  const startNexButton = () => {
    switch(gameState){
      case State.START:
        return (<button onClick={handleOnPlay} className="btn btn-lg btn-wide border-0 shadow-md text-lg text-accent bg-secondary">Start</button>)
      case State.PLAY:
        return (<div className="h-12 w-full"></div>)
      case State.ANSWERED:
        return (<button onClick={() => fetchPair()} className="btn btn-lg btn-wide border-0 shadow-md text-lg animate-in fade-in text-accent bg-secondary hover:bg-neutral">Next</button>)
      case State.END:
        return (<div className="h-12 w-full"></div>)
    }
  }

  return (
    <div className="h-screen">
      <NavBar handleChangeTheme={handleChangeTheme}></NavBar>
      <div className="h-5/6 flex flex-col justify-center items-center bg-primary text-accent">
        <div className="w-full h-1/3 flex justify-center items-center">
        {mainContent()}
        </div>
        <div className="w-full h-1/6 flex justify-center">
        {startNexButton()}
        </div>
      </div>
    </div>
  )
}

export default App
