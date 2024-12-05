
type DisplayMessageProps = {
  handleClick?: React.MouseEventHandler<HTMLButtonElement>
  headingText?: string
  buttonText?: string
}

export function DisplayMessage(props: DisplayMessageProps){
  return (
    <div className="flex font-mono text-2xl font-bold">
      <h1>{props.headingText}</h1>
      <button onClick={props.handleClick} className="underline ml-2">{props.buttonText}</button>
    </div>
  )
}