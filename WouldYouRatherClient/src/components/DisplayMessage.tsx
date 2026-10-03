
type DisplayMessageProps = {
  handleClick?: React.MouseEventHandler<HTMLButtonElement>
  headingText?: string
  buttonText?: string
}

export function DisplayMessage(props: DisplayMessageProps){
  return (
    <div className="flex flex-wrap justify-center gap-x-2 text-center text-balance font-mono text-lg sm:text-xl lg:text-2xl font-bold">
      <h1>{props.headingText}</h1>
      {props.buttonText ? <button onClick={props.handleClick} className="underline">{props.buttonText}</button> : <></>}
    </div>
  )
}