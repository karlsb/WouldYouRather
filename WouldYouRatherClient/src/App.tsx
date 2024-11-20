import { useState } from "react"



//const url = "https://setuptest-442308.lm.r.appspot.com/random-pair"
const devurl = "http://localhost:8080/random-pair"
const devurlPost = "http://localhost:8080/store-answer"

type Pair = {
  left: string
  right: string
}

type CardProps = {
  children: React.ReactNode
}

function Card(props: CardProps){
  const handleClick = async (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    e.preventDefault()
    const res = await fetch(devurlPost, {method:"POST", headers: {"Content-Type":"application/json"}}) //need to attach payload
    if(res.ok) {
      const data = await res.json()
      console.log(data)
    }
    else{
      console.log("API call failed", res.status)
    }
    

  }

  return (
    <>
      <div onClick = {(e) => handleClick(e)}className="flex justify-center items-center m-0 auto w-1/2 p-6 bg-white border border-gray-200 rounded-lg shadow hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700">
        <h2 className="font-bold text-2xl text-gray-700 dark:text-gray-400">{props.children}</h2>
      </div>
    </>
  )  
}


function App() {
  const [pair,setPair] = useState<Pair>({left:"Welcome to",right:"Would you rather"})

  async function apiCall(){
    const res = await fetch(devurl, {method:"GET", headers: {"Content-Type":"application/json"}})
    if(res.ok) {
      const data = await res.json()
      setPair({left:data.pair.left, right:data.pair.right})
    }
    else{
      console.log("API call failed", res.status)
    }
  }

  return (
    <div className="h-screen bg-slate-100">
      <div className="h-1/6 bg-blue-400 flex justify-center items-center"> {/* NavBar */}
        <h1 className='text-3xl flex font-bold underline'>
          Would You Rather, Programmer edition
        </h1>
      </div>
      <div className="h-5/6 bg-blue-100 flex flex-col justify-center items-center">{/* Main content*/}
        <div className="w-full h-5/6 flex justify-center items-center">
          <div className="w-3/5 h-4/5 flex">
            <Card>{pair.left}</Card>
            <Card>{pair.right}</Card>
          </div>
        </div>
        <div className="bg-slate-100 pl-3 pr-3 pt-1 pb-1 rounded-md shadow-md" >
          <button className="text-2xl font-bold"onClick={() => apiCall()}>Next</button>
        </div>
      </div>

    </div>
  )
}

export default App
