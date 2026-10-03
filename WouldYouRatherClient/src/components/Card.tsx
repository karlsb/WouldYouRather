import { useEffect, useState } from "react"
import { CardState } from "../types/State"

type CardProps = {
  text : string
  id: number
  choiceMade:boolean
  state: CardState
  handleClick: (e:React.MouseEvent<HTMLDivElement, MouseEvent>) => void
}

const baseClasses = "flex flex-1 justify-center items-center min-h-[13.889dvh] rounded-full px-8 py-4 sm:p-6 transition-colors duration-400"

const stateClasses: Record<CardState, string> = {
  [CardState.ShowQuestion]: "bg-secondary hover:bg-neutral hover:text-primary",
  [CardState.ShowAnswer]: "bg-secondary",
  [CardState.Picked]: "bg-neutral text-primary",
}

export function Card(props: CardProps){
  const [text, setText] = useState("")

  const onPress = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    e.preventDefault()
    if(!props.choiceMade){
        props.handleClick(e)
    }
  }

  useEffect(() => {
    setText(props.text)
  },[props])


  return (
    <>
      <div onClick = {onPress} className={`${baseClasses} ${stateClasses[props.state]}`}>
        <h2 key={text} className="text-center text-balance font-mono font-bold text-lg sm:text-xl lg:text-2xl animate-fade text-light">{text}</h2>
      </div>
    </>
  )  
}