import { useEffect, useState } from "react"

type Pair = {
  id: number
  left: string
  right: string
}

type CardProps = {
  text : string
  id: number
  side: string
  percent: number
  showPercent: boolean
  handleClick: (e:React.MouseEvent<HTMLDivElement, MouseEvent>) => void
}




const random_pair_url = import.meta.env.VITE_RANDOM_PAIR_URL
const store_answer_url = import.meta.env.VITE_STORE_ANSWER_URL

function Card(props: CardProps){
  const [text, setText] = useState("")

  let classes = props.side === "left" ?
      "border-r-2 flex justify-center items-center m-0 auto w-1/2 p-6 hover:bg-gray-100 "
    : "flex justify-center items-center m-0 auto w-1/2 p-6 hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700"

  useEffect(() => {
    setText(props.text)
  },[props])
 
  
  return (
    <>
      <div onClick = {props.handleClick} className={classes}>
        <h2 className="font-mono font-bold text-2xl text-gray-700 dark:text-gray-400">{text}</h2>
      </div>
    </>
  )  
}


type CardWrapperProps = {
  pair: Pair
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
        console.log("API call failed", res.status)
      }
    }
  }


  return (
          <div className="w-3/5 h-4/5 flex">
            <Card handleClick={(e) => handleClick("left", e)} side="left" percent={leftPercent} showPercent={choiceMade} text={leftText} id={props.pair.id}></Card>
            <Card handleClick={(e) => handleClick("right", e)} side="right" percent={rightPercent} showPercent={choiceMade} text={rightText} id={props.pair.id}></Card>
          </div>
  )
}


function App() {
  const [pair,setPair] = useState<Pair>({id:-1, left:"Welcome to",right:"Would you rather"})

  async function apiCall(){
    const res = await fetch(random_pair_url, {method:"GET", credentials:"include",headers: {"Content-Type":"application/json"}})
    console.log(res)
    if(res.ok) {
      console.log(document.cookie)
      const data = await res.json()
      setPair({id:data.pair.id, left:data.pair.left, right:data.pair.right})
    }
    else{
      console.log("API call failed", res.status)
    }
  }

  useEffect(() => {
    apiCall()
  },[])

  return (
    <div className="h-screen">
      <div className="h-1/6 flex justify-center items-center"> {/* NavBar */}
        <h1 className='text-3xl flex font-mono font-bold underline'>
          Would You Rather - Programmer edition
        </h1>
      </div>
      <div className="h-5/6 bg-base-100 flex flex-col justify-center items-center">{/* Main content*/}
        <div className="w-full h-5/6 flex justify-center items-center">
          <CardWrapper pair={pair}></CardWrapper>
        </div>
        <div className="" >
          <button className="btn bg-white border-0 shadow-none text-lg"onClick={() => apiCall()}>Next</button>
        </div>
      </div>
    </div>
  )
}

export default App
