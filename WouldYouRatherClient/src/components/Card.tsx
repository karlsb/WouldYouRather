import { useEffect, useState } from "react"

type CardProps = {
  text : string
  id: number
  side: string
  percent: number
  showPercent: boolean
  handleClick: (e:React.MouseEvent<HTMLDivElement, MouseEvent>) => void
}

export function Card(props: CardProps){
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