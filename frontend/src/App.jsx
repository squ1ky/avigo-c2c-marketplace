import { BrowserRouter as Router } from "react-router-dom"
import { Toaster } from 'react-hot-toast'

function App() {
    return (
        <Router>
            <Toaster position="top-right" />
            <div>
                <h1>AviGo Frontend - Ready!</h1>
            </div>
        </Router>
    )
}