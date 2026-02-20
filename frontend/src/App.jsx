import { use, useState } from "react"

function App() {

  let [text , setText] = useState("")
  let [audioUrl , setAudioUrl] = useState("")
  const [selectedFile, setSelectedFile] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const handleConvert = async () => {
    const response = await fetch('http://localhost:8080/createRequest',{
      method: "POST",
      body: JSON.stringify({ text: text }),
      headers: { "Content-Type": "application/json" },
    })
    const data = await response.json();
    setaudioUrl("http://localhost:8080" + data.audio_url);
  };

  const handleFileChange = (event) => {
    // Save the file object from the input
    if (event.target.files && event.target.files[0]) {
      setSelectedFile(event.target.files[0]);
    }
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      alert("Please select a file first!");
      return;
    }

    setIsLoading(true);

    try {
      // Create a "FormData" package (just like Postman form-data)
      const formData = new FormData();
      formData.append("file", selectedFile); // Key must be "file" to match c.FormFile("file")

      // Send POST request
      // NOTE: Make sure this URL matches your Go route ("/upload" or "/uploadText")
      const response = await fetch("http://localhost:8080/upload", {
        method: "POST",
        body: formData,
      });

      if (!response.ok) {
        throw new Error("Upload failed");
      }

      const data = await response.json();
      
      // 3. Set the Audio URL to play it
      // Ensure the backend returns { "audio_url": "/audio/..." }
      setAudioUrl("http://localhost:8080" + data.audio_url);

    } catch (error) {
      console.error("Error uploading file:", error);
      alert("Something went wrong! Check console.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
  <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-center p-6">
    <div className="max-w-md w-full bg-white rounded-2xl shadow-xl p-8 border border-gray-100">
      <h1 className="text-2xl font-bold text-gray-800 text-center mb-2">
        Speech Generator
      </h1>
      <p className="text-gray-500 text-center mb-8 text-sm">
        Upload a .txt file to convert it into high-quality audio
      </p>

      <div className="space-y-6">
        {/* --- CUSTOM FILE INPUT --- */}
        <div className="relative">
          <label 
            className={`flex flex-col items-center justify-center w-full h-32 border-2 border-dashed rounded-xl cursor-pointer transition-all
              ${selectedFile ? 'border-green-400 bg-green-50' : 'border-gray-300 bg-gray-50 hover:bg-gray-100'}`}
          >
            <div className="flex flex-col items-center justify-center pt-5 pb-6">
              <svg className="w-8 h-8 mb-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
              </svg>
              <p className="text-sm text-gray-600">
                {selectedFile ? (
                  <span className="font-semibold text-green-600">{selectedFile.name}</span>
                ) : (
                  <span>Click to upload or drag and drop</span>
                )}
              </p>
            </div>
            {/* The actual hidden input */}
            <input 
              type="file" 
              className="hidden" 
              onChange={handleFileChange} 
            />
          </label>
        </div>

        {/* --- ACTION BUTTON --- */}
        <button
          onClick={handleUpload}
          disabled={isLoading || !selectedFile}
          className={`w-full py-3 px-4 rounded-xl font-semibold text-white transition-all shadow-md active:scale-95
            ${isLoading || !selectedFile 
              ? 'bg-gray-400 cursor-not-allowed' 
              : 'bg-indigo-600 hover:bg-indigo-700 hover:shadow-indigo-200'}`}
        >
          {isLoading ? (
            <span className="flex items-center justify-center">
              <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Processing Audio...
            </span>
          ) : "Generate Speech"}
        </button>

        {/* --- AUDIO PLAYER SECTION --- */}
        {audioUrl && (
          <div className="mt-8 p-4 bg-indigo-50 rounded-xl border border-indigo-100 animate-fade-in">
            <h3 className="text-sm font-bold text-indigo-900 mb-3 flex items-center">
              <span className="mr-2">🎧</span> Resulting Audio
            </h3>
            <audio controls src={audioUrl} className="w-full h-10">
              Your browser does not support the audio element.
            </audio>
          </div>
        )}
      </div>
    </div>
  </div>
);

}

export default App
