
import { useEffect, useState } from "react"
import { Card } from "./Card"
import Pair from "../types/Pair"

type CardWrapperProps = {
  pair: Pair
  handleAnswer?: (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => void
}

export function CardWrapper(props: CardWrapperProps) {
  const [rightText, setRightText] = useState("")
  const [leftText, setLeftText] = useState("")
  const [leftPercent, setLeftPercent] = useState(0)
  const [rightPercent, setRightPercent] = useState(0)
  const [choiceMade, setChoiceMade] = useState(false)
    
  const store_answer_url = import.meta.env.VITE_STORE_ANSWER_URL

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
    }
  }

  return (
          <div className="w-3/5 flex flex-wrap animate-in slide-in-from-left bg-primary">
            <Card handleClick={(e) => handleClick("left", e)} side="left" percent={leftPercent} showPercent={choiceMade} text={leftText} id={props.pair.id}></Card>
            <Card handleClick={(e) => handleClick("right", e)} side="right" percent={rightPercent} showPercent={choiceMade} text={rightText} id={props.pair.id}></Card>
          </div>
  )
}
