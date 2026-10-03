
import { useEffect, useState } from "react"
import { Card } from "./Card"
import Pair from "../types/Pair"
import { CardState, GameState } from "../types/State"

type CardWrapperProps = {
  pair: Pair
  handleAnswer?: (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => void
  state: GameState
}

export function CardWrapper(props: CardWrapperProps) {
  const [rightText, setRightText] = useState("")
  const [leftText, setLeftText] = useState("")
  const [leftPercent, setLeftPercent] = useState("")
  const [rightPercent, setRightPercent] = useState("")
  const [choiceMade, setChoiceMade] = useState(false)

  const [leftCardState, setLeftCardState] = useState(CardState.ShowQuestion)
  const [rightCardState, setRightCardState] = useState(CardState.ShowQuestion)
    
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
        setLeftPercent(lpercentage.toFixed(0).toString() + "% Picked this option")
        setRightPercent(rpercentage.toFixed(0).toString()+ "% Picked this option")
        setChoiceMade(true)
        if(leftright === "left"){
          setLeftCardState(CardState.Picked)
          setRightCardState(CardState.ShowAnswer)
        }else if (leftright === "right"){
          setLeftCardState(CardState.ShowAnswer)
          setRightCardState(CardState.Picked)
        }
      }
    }
  }

  useEffect(() => {
    if(props.state === GameState.PLAY){
      setLeftCardState(CardState.ShowQuestion)
      setRightCardState(CardState.ShowQuestion)
      setChoiceMade(false)
    }
  },[props.state])


  return (
          <div className="self-start w-full lg:w-4/5 max-w-6xl grid auto-rows-fr sm:grid-cols-2 gap-x-4 animate-in slide-in-from-left bg-primary">
            <div className="flex flex-col">
                <Card handleClick={(e) => handleClick("left", e)} state={leftCardState} choiceMade={choiceMade} text={leftText} id={props.pair.id}></Card>
                <div className="min-h-7 lg:min-h-8 mt-2 mb-5 sm:mt-5 sm:mb-4">
                  {choiceMade ? <p className="text-center animate-fade font-mono font-bold text-lg sm:text-xl lg:text-2xl">{leftPercent}</p> : <></>}
                </div>
            </div>
            <div className="flex flex-col">
                <Card handleClick={(e) => handleClick("right", e)} state={rightCardState} choiceMade={choiceMade} text={rightText} id={props.pair.id}></Card>
                <div className="min-h-7 lg:min-h-8 mt-2 mb-5 sm:mt-5 sm:mb-4">
                  {choiceMade ? <p className="text-center animate-fade font-mono font-bold text-lg sm:text-xl lg:text-2xl">{rightPercent}</p> : <></>}
                </div>
            </div>
          </div>
  )
}
