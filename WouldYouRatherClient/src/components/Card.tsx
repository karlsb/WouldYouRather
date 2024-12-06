import { useEffect, useState } from "react"
import { CardState } from "../types/State"

type CardProps = {
  text : string
  id: number
  side: string
  choiceMade:boolean
  state: CardState
  handleClick: (e:React.MouseEvent<HTMLDivElement, MouseEvent>) => void
}

export function Card(props: CardProps){
  const [text, setText] = useState("")
  const [additional, setAdditional] = 
     useState( props.side === "left" ? 
        "flex justify-center items-center h-1/2 rounded-full mr-2 m-0 auto p-6 bg-secondary transition-colors duration-400 mr-2"
       :"flex justify-center items-center h-1/2 rounded-full ml-2 m-0 auto p-6 bg-secondary transition-colors duration-400 ml-2"
    )

  //let classes = "flex justify-center items-center h-1/2 rounded-full mr-2 m-0 auto p-6 bg-secondary transition-colors duration-400"
  //classes = props.side === "left" ?
        //classes + " mr-2"
      //: classes + " ml-2"

  //classes = props.choiceMade ? classes : classes + " hover:bg-neutral"

  const onPress = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    e.preventDefault()
    if(!props.choiceMade){
        props.handleClick(e)
    }
  }

  useEffect(() => {
    setText(props.text)
    if(props.state === CardState.Picked){
      setAdditional( props.side === "left" ?  
        "flex justify-center items-center h-1/2 rounded-full mr-2 m-0 auto p-6 transition-colors duration-400 mr-2 bg-neutral text-primary"
       :"flex justify-center items-center h-1/2 rounded-full mr-2 m-0 auto p-6 transition-colors duration-400 ml-2 bg-neutral text-primary"
      )
    }
    else if (props.state === CardState.ShowAnswer){
      setAdditional( props.side === "left" ?  
        "flex justify-center items-center h-1/2 rounded-full mr-2 m-0 auto p-6 transition-colors duration-400 mr-2 bg-secondary"
       :"flex justify-center items-center h-1/2 rounded-full mr-2 m-0 auto p-6 transition-colors duration-400 ml-2 bg-secondary"
      )
    }
    else if (props.state === CardState.ShowQuestion){
      setAdditional( props.side === "left" ?  
        "flex justify-center items-center h-1/2 rounded-full mr-2 m-0 auto p-6 transition-colors duration-400 mr-2 bg-secondary hover:bg-neutral hover:text-primary"
       :"flex justify-center items-center h-1/2 rounded-full mr-2 m-0 auto p-6 transition-colors duration-400 ml-2 bg-secondary hover:bg-neutral hover:text-primary"
      )
    }
  },[props])


  return (
    <>
      <div onClick = {onPress} className={additional}>
        <h2 key={text} className="text-center text-balance font-mono font-bold text-2xl animate-fade text-light">{text}</h2>
      </div>
    </>
  )  
}