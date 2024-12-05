type NavBarProps = {
  handleChangeTheme: (newTheme: string) => void
}

export function NavBar(props: NavBarProps) {

  return (
      <div className="h-1/6 flex justify-center items-center relative bg-secondary px-4" >
        <h1 className="text-3xl font-extrabold animate-fade text-accent mx-auto">
          Would You Rather - Programmer Edition
        </h1>
      <div className="absolute right-40">
        <details className="dropdown">
          <summary className="btn m-1">Color Theme</summary>
          <ul className="menu dropdown-content bg-base-100 rounded-box z-[1] w-52 p-2 shadow">
            <li><button onClick={() => props.handleChangeTheme('one')}>Color Theme 1</button></li>
            <li><button onClick={() => props.handleChangeTheme('two')}>Color Theme 2</button></li>
            <li><button onClick={() => props.handleChangeTheme('three')}>Color Theme 3</button></li>
            <li><button onClick={() => props.handleChangeTheme('four')}>Color Theme 4</button></li>
            <li><button onClick={() => props.handleChangeTheme('five')}>Color Theme 5</button></li>
            <li><button onClick={() => props.handleChangeTheme('six')}>Color Theme 6</button></li>
            <li><button onClick={() => props.handleChangeTheme('seven')}>Color Theme 7</button></li>
          </ul>
        </details>
        </div>
      </div> 
        )
}