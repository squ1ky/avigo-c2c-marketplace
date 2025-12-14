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
import AccountPage from "./pages/AccountPage.jsx";
import EditProfilePage from "./pages/EditProfilePage.jsx";
import ChangePasswordPage from "./pages/ChangePasswordPage.jsx";
import CreateListingPage from "./pages/CreateListingPage.jsx";
import ListingPage from "./pages/ListingPage.jsx";
import ProtectedRoute from "./components/auth/ProtectedRoute.jsx";
import MySalesPage from "./pages/MySalesPage.jsx";
import MyPurchasesPage from "./pages/MyPurchasesPage.jsx";
import MyListingsPage from "./pages/MyListingsPage.jsx";

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
                    <Route path="/listing/:id" element={<Layout><ListingPage /></Layout>} />

                    {/* protected routes */}
                    <Route path="/create" element={<ProtectedRoute><Layout><CreateListingPage /></Layout></ProtectedRoute>}/>
                    <Route path="/account" element={<ProtectedRoute><Layout><AccountPage /></Layout></ProtectedRoute>}/>
                    <Route path="/account/edit" element={<ProtectedRoute><Layout><EditProfilePage /></Layout></ProtectedRoute>}/>
                    <Route path="/account/change-password" element={<ProtectedRoute><Layout><ChangePasswordPage /></Layout></ProtectedRoute>}/>
                    <Route path="/account/purchases" element={<ProtectedRoute><Layout><MyPurchasesPage /></Layout></ProtectedRoute>}/>
                    <Route path="/account/sales" element={<ProtectedRoute><Layout><MySalesPage /></Layout></ProtectedRoute>}/>
                    <Route path="/account/listings" element={<ProtectedRoute><Layout><MyListingsPage /></Layout></ProtectedRoute>}/>
                </Routes>
            </AuthProvider>
        </Router>
    );
}

export default App