type NavBarProps = {
  handleChangeTheme: (newTheme: string) => void
}

export function NavBar(props: NavBarProps) {

  return (
      <div className="min-h-[16.667dvh] flex items-center gap-3 relative bg-secondary px-4 py-2" >
        <h1 className="flex-1 text-lg sm:text-xl md:text-2xl lg:text-3xl lg:text-center text-balance font-extrabold animate-fade text-accent">
          Would You Rather - Programmer Edition
        </h1>
      <div className="lg:absolute lg:right-4 xl:right-20 min-[1440px]:right-40">
        <details className="dropdown dropdown-end">
          <summary className="btn m-1 bg-primary border-none shadow-md"><span><span className="hidden sm:inline">Color </span>Theme</span></summary>
          <ul className="menu dropdown-content bg-base-100 rounded-box z-[1] w-52 p-2 shadow">
            <li><button className="hover:bg-primary" onClick={() => props.handleChangeTheme('one')}>Color Theme 1</button></li>
            <li><button className="hover:bg-primary" onClick={() => props.handleChangeTheme('two')}>Color Theme 2</button></li>
            <li><button className="hover:bg-primary" onClick={() => props.handleChangeTheme('three')}>Color Theme 3</button></li>
            <li><button className="hover:bg-primary" onClick={() => props.handleChangeTheme('four')}>Color Theme 4</button></li>
            <li><button className="hover:bg-primary" onClick={() => props.handleChangeTheme('five')}>Color Theme 5</button></li>
            <li><button className="hover:bg-primary" onClick={() => props.handleChangeTheme('six')}>Color Theme 6</button></li>
            <li><button className="hover:bg-primary" onClick={() => props.handleChangeTheme('seven')}>Color Theme 7</button></li>
          </ul>
        </details>
        </div>
      </div> 
        )
}