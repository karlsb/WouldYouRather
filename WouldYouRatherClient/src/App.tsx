import { useEffect, useState } from "react"


// colors from : https://colorhunt.co/palette/89a8b2b3c8cfe5e1daf1f0e8

type Pair = {
  id: number
  left: string
  right: string
}

const random_pair_url = import.meta.env.VITE_RANDOM_PAIR_URL
const store_answer_url = import.meta.env.VITE_STORE_ANSWER_URL
const n_random_pairs_url = import.meta.env.VITE_N_RANDOM_PAIRS_URL

type DisplayMessageProps = {
  handleClick?: React.MouseEventHandler<HTMLButtonElement>
  headingText?: string
  buttonText?: string
}

function DisplayMessage(props: DisplayMessageProps){
  return (
    <div className="flex font-mono text-2xl font-bold">
      <h1>{props.headingText}</h1>
      <button onClick={props.handleClick} className="underline ml-2">{props.buttonText}</button>
    </div>
  )
}


type CardProps = {
  text : string
  id: number
  side: string
  percent: number
  showPercent: boolean
  handleClick: (e:React.MouseEvent<HTMLDivElement, MouseEvent>) => void
}

function Card(props: CardProps){
  const [text, setText] = useState("")

  let classes = props.side === "left" ?
        "flex flex-1 flex-col justify-center rounded-full mr-2 items-center m-0 auto w-1/2 p-6 bg-secondary hover:bg-neutral transition-colors duration-400"
      : "flex flex-1 flex-col justify-center rounded-full ml-2 m-0 auto w-1/2 p-6 bg-secondary hover:bg-neutral transition-colors duration-400"

  useEffect(() => {
    setText(props.text)
  },[props])
 
  return (
    <>
      <div onClick = {props.handleClick} className={classes}>
        <h2 key={text} className="text-center text-balance font-mono font-bold text-2xl animate-fade text-light" >{text}</h2>
      </div>
    </>
  )  
}

type CardWrapperProps = {
  pair: Pair
  handleAnswer?: (e:React.MouseEvent<HTMLDivElement, MouseEvent>) => void
}

function CardWrapper(props: CardWrapperProps) {
  const [rightText, setRightText] = useState("")
  const [leftText, setLeftText] = useState("")
  const [leftPercent, setLeftPercent] = useState(0)
  const [rightPercent, setRightPercent] = useState(0)
  const [choiceMade, setChoiceMade] = useState(false)

  useEffect(() => {
    setLeftText(props.pair.left)
    setRightText(props.pair.right)
    setChoiceMade(false)
  },[props])

  const handleClick = async (leftright: string, e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    e.preventDefault()
    if(props.handleAnswer){
      props.handleAnswer(e)
    }
    if(!choiceMade) {
      const res = await fetch(store_answer_url, {method:"POST", credentials:"include", headers: {"Content-Type":"application/json"}, body: JSON.stringify({"id":props.pair.id , "leftright": leftright })}) //need to attach payload
      if(res.ok) {
        const data = await res.json()
        const lpercentage = data.pair.lcount * 100 / (data.pair.lcount + data.pair.rcount) 
        const rpercentage = data.pair.rcount * 100 / (data.pair.lcount + data.pair.rcount)  
        setLeftText(lpercentage.toFixed(0).toString() + "% Picked this option")
        setLeftPercent(lpercentage)
        setRightText(rpercentage.toFixed(0).toString()+ "% Picked this option")
        setRightPercent(rpercentage)
        setChoiceMade(true)
      }
      else{
      }
    }
  }

  return (
          <div className="w-3/5 flex flex-wrap animate-in slide-in-from-left bg-primary">
            <Card handleClick={(e) => handleClick("left", e)} side="left" percent={leftPercent} showPercent={choiceMade} text={leftText} id={props.pair.id}></Card>
            <Card handleClick={(e) => handleClick("right", e)} side="right" percent={rightPercent} showPercent={choiceMade} text={rightText} id={props.pair.id}></Card>
          </div>
  )
}

type NavBarProps = {
  handleChangeTheme: (newTheme: string) => void
}
function NavBar(props: NavBarProps) {

  return (
      <div className="h-1/6 flex justify-center items-center relative bg-secondary px-4" >
        <h1 className="text-3xl font-extrabold animate-fade text-accent mx-auto">
          Would You Rather - Programmer Edition
        </h1>
      <div className="absolute right-40">
        <details className="dropdown">
          <summary className="btn m-1">Color Theme</summary>
          <ul className="menu dropdown-content bg-base-100 rounded-box z-[1] w-52 p-2 shadow">
            <li><button onClick={() => props.handleChangeTheme('one')}>Color Theme 1</button></li>
            <li><button onClick={() => props.handleChangeTheme('two')}>Color Theme 2</button></li>
            <li><button onClick={() => props.handleChangeTheme('three')}>Color Theme 3</button></li>
            <li><button onClick={() => props.handleChangeTheme('four')}>Color Theme 4</button></li>
            <li><button onClick={() => props.handleChangeTheme('five')}>Color Theme 5</button></li>
            <li><button onClick={() => props.handleChangeTheme('six')}>Color Theme 6</button></li>
            <li><button onClick={() => props.handleChangeTheme('seven')}>Color Theme 7</button></li>
          </ul>
        </details>
        </div>
      </div> 
        )
}

function App() {
  const themes = ["one","two", "three", "four","five","six","seven"]
  const [theme, setTheme] = useState("one")
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
      setTheme(savedTheme)
      document.documentElement.setAttribute('data-theme', savedTheme)
      console.log(theme)
      console.log(themes)
    }
  },[])

  function handleChangeTheme(newTheme:string) {
    console.log("CHANGING COLOR THEME", newTheme)
    setTheme(newTheme)
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

  //TODO make a render switch statement with 3 states - start, playing, endstate - render this function the div
  //TODO where the allPairsSeen is being rendered now

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
        return (<button onClick={handleOnPlay} className="btn border-0 shadow-none text-lg text-accent bg-secondary">Start</button>)
      case State.PLAY:
        return (<div className="h-12 w-full"></div>)
      case State.ANSWERED:
        return (<button onClick={() => fetchPair()} className="btn border-0 shadow-none text-lg animate-in fade-in text-accent bg-secondary hover:bg-neutral">Next</button>)
      case State.END:
        return (<div className="h-12 w-full"></div>)
    }
  }

  return (
    <div className="h-screen">
      <NavBar handleChangeTheme={handleChangeTheme}></NavBar>
      <div className="h-5/6 flex flex-col justify-center items-center bg-primary text-accent">
        <div className="w-full h-5/6 flex justify-center items-center">
        {mainContent()}
        </div>
        {startNexButton()}
      </div>
    </div>
  )
}

export default App
