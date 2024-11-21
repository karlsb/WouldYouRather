import { useEffect, useState } from "react"

type Pair = {
  id: number
  left: string
  right: string
}

type CardProps = {
  text : string
  id: number
  handleClick: (e:React.MouseEvent<HTMLDivElement, MouseEvent>) => void
}

const random_pair_url = import.meta.env.VITE_RANDOM_PAIR_URL
const store_answer_url = import.meta.env.VITE_STORE_ANSWER_URL

//TODO consider sending a full pair as props, to make it cleaner
function Card(props: CardProps){
  const [text, setText] = useState("")

  useEffect(() => {
    setText(props.text)
  },[props])
  
  return (
    <>
      <div onClick = {props.handleClick} className="flex justify-center items-center m-0 auto w-1/2 p-6 bg-white border border-gray-200 rounded-lg shadow hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700">
        <h2 className="font-bold text-2xl text-gray-700 dark:text-gray-400">{text}</h2>
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
  const [choiceMade, setChoiceMade] = useState(false)

  useEffect(() => {
    setLeftText(props.pair.left)
    setRightText(props.pair.right)
    setChoiceMade(false)
  },[props])

  const handleClick = async (leftright: string, e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    e.preventDefault()
    if(!choiceMade) {
      const res = await fetch(store_answer_url, {method:"POST", headers: {"Content-Type":"application/json"}, body: JSON.stringify({"id":props.pair.id , "leftright": leftright })}) //need to attach payload
      if(res.ok) {
        const data = await res.json()
        const lpercentage = data.pair.lcount * 100 / (data.pair.lcount + data.pair.rcount) 
        const rpercentage = data.pair.rcount * 100 / (data.pair.lcount + data.pair.rcount)  
        setLeftText(leftText + " " + lpercentage.toFixed(0).toString() + "% Picked this option")
        setRightText(rightText + " " + rpercentage.toFixed(0).toString()+ "% Picked this option")
        setChoiceMade(true)
      }
      else{
        console.log("API call failed", res.status)
      }
    }
  }


  return (
          <div className="w-3/5 h-4/5 flex">
            <Card handleClick={(e) => handleClick("left", e)} text={leftText} id={props.pair.id}></Card>
            <Card handleClick={(e) => handleClick("right", e)} text={rightText} id={props.pair.id}></Card>
          </div>
  )
}


function App() {
  const [pair,setPair] = useState<Pair>({id:-1, left:"Welcome to",right:"Would you rather"})

  async function apiCall(){
    const res = await fetch(random_pair_url, {method:"GET", headers: {"Content-Type":"application/json"}})
    if(res.ok) {
      const data = await res.json()
      setPair({id:data.pair.id, left:data.pair.left, right:data.pair.right})
    }
    else{
      console.log("API call failed", res.status)
    }
  }

  return (
    <div className="h-screen bg-slate-100">
      <div className="h-1/6 bg-blue-400 flex justify-center items-center"> {/* NavBar */}
        <h1 className='text-3xl flex font-bold underline'>
          Would You Rather, Programmer edition
        </h1>
      </div>
      <div className="h-5/6 bg-blue-100 flex flex-col justify-center items-center">{/* Main content*/}
        <div className="w-full h-5/6 flex justify-center items-center">
          <CardWrapper pair={pair}></CardWrapper>
        </div>
        <div className="" >
          <button className="btn"onClick={() => apiCall()}>Next</button>
        </div>
      </div>
    </div>
  )
}

export default App
