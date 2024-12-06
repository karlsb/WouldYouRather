import { useEffect, useState } from "react"
import className from 'classnames'

type CardProps = {
  text : string
  id: number
  side: string
  choiceMade:boolean
  handleClick: (e:React.MouseEvent<HTMLDivElement, MouseEvent>) => void
}

export function Card(props: CardProps){
  const [text, setText] = useState("")
  const [additional, setAdditional] = 
     useState( props.side === "left" ? 
        "flex justify-center items-center h-1/2 rounded-full mr-2 m-0 auto p-6 bg-secondary transition-colors duration-400 mr-2"
       :"flex justify-center items-center h-1/2 rounded-full mr-2 m-0 auto p-6 bg-secondary transition-colors duration-400 ml-2"

    )
  const [picked, setPicked] = useState(false)

  let classes = "flex justify-center items-center h-1/2 rounded-full mr-2 m-0 auto p-6 bg-secondary transition-colors duration-400"
  classes = props.side === "left" ?
        classes + " mr-2"
      : classes + " ml-2"

  classes = props.choiceMade ? classes : classes + " hover:bg-neutral"

  const onPress = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    e.preventDefault()
    setPicked(true)
    if(!props.choiceMade){
        props.handleClick(e)
    }
  }

  useEffect(() => {
    if(picked){
    setAdditional((prev) => {
      const updatedClasses = prev
        .split(" ")
        .filter((cls) => cls !== "bg-secondary") // Remove the target class
        .join(" ");
      return `${updatedClasses}  bg-neutral text-primary`; // Add the new class
    })
    }
  },[picked])

  useEffect(() => {
    if(!props.choiceMade){
        setAdditional("flex justify-center items-center h-1/2 rounded-full mr-2 m-0 auto p-6 bg-secondary transition-colors duration-400 mr-2")
    }
  },[props.choiceMade])

  useEffect(() => {
    setText(props.text)
    setPicked(false)
  },[props])


  return (
    <>
      <div onClick = {onPress} className={additional}>
        <h2 key={text} className="text-center text-balance font-mono font-bold text-2xl animate-fade text-light">{text}</h2>
      </div>
    </>
  )  
}