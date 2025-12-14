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

                    <Route element={<ProtectedRoute />}>
                        <Route path="/create" element={<Layout><CreateListingPage /></Layout>} />

                        <Route path="/account" element={<Layout><AccountPage /></Layout>} />
                        <Route path="/account/edit" element={<Layout><EditProfilePage /></Layout>} />
                        <Route path="/account/change-password" element={<Layout><ChangePasswordPage /></Layout>} />

                        <Route path="/account/purchases" element={<Layout><MyPurchasesPage /></Layout>} />
                        <Route path="/account/sales" element={<Layout><MySalesPage /></Layout>} />
                        <Route path="/account/listings" element={<Layout><MyListingsPage /></Layout>} />
                    </Route>
                </Routes>
            </AuthProvider>
        </Router>
    );
}

export default App