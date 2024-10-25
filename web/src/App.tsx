import Content from "./components/Content"
import WebSocketContextProvider from "./contexts/WebSocketContext"

function App() {
  return (
    <WebSocketContextProvider>
      <div className="w-screen h-screen bg-slate-800 flex justify-center items-center">
        <Content />
      </div>
    </WebSocketContextProvider>
  )
}

export default App
