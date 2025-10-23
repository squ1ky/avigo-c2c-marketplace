import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import { AuthProvider } from './context/AuthContext'
import Layout from './components/layout/Layout'
import HomePage from './pages/HomePage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import EmailConfirmPage from './pages/EmailConfirmPage'

import './styles/header.css'
import './styles/footer.css'

function App() {
    return (
        <Router>
            <AuthProvider>
                <Toaster position="top-right" />
                <Routes>
                    <Route path="/" element={<Layout><HomePage /></Layout>} />
                    <Route path="/auth/login" element={<Layout><LoginPage /></Layout>} />
                    <Route path="/auth/register" element={<Layout><RegisterPage /></Layout>} />
                    <Route path="/auth/confirm-email" element={<Layout><EmailConfirmPage /></Layout>} />
                </Routes>
            </AuthProvider>
        </Router>
    )
}

export default App