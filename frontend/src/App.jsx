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
import UserSalesPage from "./pages/UserSalesPage.jsx";
import UserPurchasesPage from "./pages/UserPurchasesPage.jsx";
import UserListingsPage from "./pages/UserListingsPage.jsx";
import EditListingPage from "./pages/EditListingPage.jsx";

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

                    <Route path="/account/:userId/listings/:listingId" element={<Layout><ListingPage /></Layout>} />

                    <Route element={<ProtectedRoute />}>
                        <Route path="/account/:userId/listings/create" element={<Layout><CreateListingPage /></Layout>} />
                        <Route path="/account/:userId/listings/:listingId/edit" element={<Layout><EditListingPage /></Layout>} />

                        <Route path="/account/:userId" element={<Layout><AccountPage /></Layout>} />
                        <Route path="/account/:userId/edit" element={<Layout><EditProfilePage /></Layout>} />
                        <Route path="/account/:userId/change-password" element={<Layout><ChangePasswordPage /></Layout>} />

                        <Route path="/account/:userId/purchases" element={<Layout><UserPurchasesPage /></Layout>} />
                        <Route path="/account/:userId/sales" element={<Layout><UserSalesPage /></Layout>} />
                        <Route path="/account/:userId/listings" element={<Layout><UserListingsPage /></Layout>} />
                    </Route>
                </Routes>
            </AuthProvider>
        </Router>
    );
}

export default App