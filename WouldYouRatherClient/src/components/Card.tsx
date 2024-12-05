import { useEffect, useState } from "react"

type CardProps = {
  text : string
  id: number
  side: string
  choiceMade:boolean
  handleClick: (e:React.MouseEvent<HTMLDivElement, MouseEvent>) => void
}

export function Card(props: CardProps){
  const [text, setText] = useState("")
  const [additional, setAdditional] = useState("")
  
  let classes = "flex justify-center items-center h-1/2 rounded-full mr-2 m-0 auto p-6 bg-secondary transition-colors duration-400"
  classes = props.side === "left" ?
        classes + " mr-2"
      : classes + " ml-2"

  classes = props.choiceMade ? classes : classes + " hover:bg-neutral"

  const onPress = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    e.preventDefault()
    if(!props.choiceMade){
        setAdditional("dummy-class bg-neutral text-primary")
        props.handleClick(e)
    }
  }

  useEffect(() => {
    if(!props.choiceMade){
        setAdditional("")
    }
  },[props.choiceMade])

  useEffect(() => {
    setText(props.text)
  },[props])
 
  return (
    <>
      <div onClick = {onPress} className={classes + additional}>
        <h2 key={text} className="text-center text-balance font-mono font-bold text-2xl animate-fade text-light" >{text}</h2>
      </div>
    </>
  )  
}