import { use, useState } from "react"

function App() {

  let [text , setText] = useState("")
  let [audioUrl , setaudioUrl] = useState("")
  const handleConvert = async () => {
    const response = await fetch('http://localhost:8080/createRequest',{
      method: "POST",
      body: JSON.stringify({ text: text }),
      headers: { "Content-Type": "application/json" },
    })
    const data = await response.json();
    setaudioUrl("http://localhost:8080" + data.audio_url);
  };

  return (
    <>
      <h1>Hello World</h1>
      
      <textarea 
      value={text}
      onChange={(e) => setText(e.target.value)}
      className="w-full p-4 mb-4 bg-white border border-gray-200 rounded-xl shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent resize-y transition-all"
      ></textarea>
      <button 
      onClick={handleConvert}
      className="w-full sm:w-auto bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-3 px-6 rounded-lg shadow-md hover:shadow-lg transition-all duration-200 mb-10"
      >Convert to Speech</button>
      
      {audioUrl && (
        <div className="p-6 bg-gray-50 border border-gray-100 rounded-2xl shadow-inner animate-fade-in-down">
          <h3 className="text-lg font-semibold mb-4 text-gray-700">
            Your Generated Audio:
          </h3>
          <audio src={audioUrl} controls autoPlay className="w-full outline-none" />
        </div>
      )}

    </>
  )

}

export default App
