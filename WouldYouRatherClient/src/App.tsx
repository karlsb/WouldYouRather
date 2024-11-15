

function Card(){
  const handleClick = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    e.preventDefault()
    console.log("We do a request to the server")
  }
  return (
    <>
      <div onClick = {(e) => handleClick(e)}className="flex justify-center items-center m-0 auto w-1/2 p-6 bg-white border border-gray-200 rounded-lg shadow hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700">
        <h2 className="font-bold text-2xl text-gray-700 dark:text-gray-400">This is a sample text</h2>
      </div>
    </>
  )  
}

function App() {

  return (
    <div className="h-screen bg-slate-100">
      <div className="h-1/6 bg-red-100 flex justify-center items-center"> {/* NavBar */}
        <h1 className='text-3xl flex font-bold underline'>
          Would You Rather !!
        </h1>
      </div>
      <div className="h-5/6 bg-blue-100 flex justify-center items-center">{/* Main content*/}
        <div className="w-3/5 h-4/5 flex">
          <Card></Card>
          <Card></Card>
        </div>
      </div>
    </div>
  )
}

export default App
