import { Link } from 'react-router-dom'
import '../styles/hero.css'
import '../styles/features.css'

function HomePage() {
    return (
        <>
            {/* Hero Section */}
            <section className="hero">
                <div className="hero-content">
                    <h1 className="hero-title">
                        <span className="hero-avi">Надежный</span> как Avito<br />
                        <span className="hero-go">Быстрый</span> как Go
                    </h1>
                    <p className="hero-subtitle">
                        Современная C2C платформа для покупки и продажи товаров
                    </p>
                    <div className="hero-actions">
                        <Link to="/auth/register" className="btn btn-large btn-primary">Начать продавать</Link>
                        <Link to="/catalog" className="btn btn-large btn-outline">Найти товар</Link>
                    </div>
                </div>
            </section>

            {/* Features Section */}
            <section className="features" id="features">
                <div className="container">
                    <h2 className="section-title">Почему AviGo?</h2>
                    <div className="features-grid">
                        <div className="feature-card">
                            <span className="material-icons feature-icon">bolt</span>
                            <h3>Молниеносная скорость</h3>
                            <p>Архитектура на Go обеспечивает мгновенную загрузку и отклик системы</p>
                        </div>
                        <div className="feature-card">
                            <span className="material-icons feature-icon">lock</span>
                            <h3>Безопасные сделки</h3>
                            <p>Проверка продавцов, защита платежей и гарантия возврата средств</p>
                        </div>
                        <div className="feature-card">
                            <span className="material-icons feature-icon">phone_android</span>
                            <h3>Удобный интерфейс</h3>
                            <p>Интуитивно понятный дизайн для комфортной работы с любого устройства</p>
                        </div>
                    </div>
                </div>
            </section>

            {/* How It Works */}
            <section className="how-it-works" id="how-it-works">
                <div className="container">
                    <h2 className="section-title">Как это работает</h2>
                    <div className="steps">
                        <div className="step">
                            <div className="step-number">1</div>
                            <h3>Зарегистрируйтесь</h3>
                            <p>Создайте аккаунт за 30 секунд</p>
                        </div>
                        <div className="step">
                            <div className="step-number">2</div>
                            <h3>Разместите объявление</h3>
                            <p>Добавьте заголовок и описание товара – наши алгоритмы сами предложат категорию и теги</p>
                        </div>
                        <div className="step">
                            <div className="step-number">3</div>
                            <h3>Продавайте</h3>
                            <p>Общайтесь с покупателями и совершайте сделки</p>
                        </div>
                    </div>
                </div>
            </section>

            {/* Categories */}
            <section className="categories" id="categories">
                <div className="container">
                    <h2 className="section-title">Популярные категории</h2>
                    <div className="categories-grid">
                        <div className="category-card">
                            <span className="material-icons category-icon">smartphone</span>
                            <div className="category-name">Электроника</div>
                        </div>
                        <div className="category-card">
                            <span className="material-icons category-icon">checkroom</span>
                            <div className="category-name">Одежда</div>
                        </div>
                        <div className="category-card">
                            <span className="material-icons category-icon">chair</span>
                            <div className="category-name">Мебель</div>
                        </div>
                        <div className="category-card">
                            <span className="material-icons category-icon">directions_car</span>
                            <div className="category-name">Автомобили</div>
                        </div>
                        <div className="category-card">
                            <span className="material-icons category-icon">home</span>
                            <div className="category-name">Недвижимость</div>
                        </div>
                        <div className="category-card">
                            <span className="material-icons category-icon">sports_soccer</span>
                            <div className="category-name">Спорт</div>
                        </div>
                    </div>
                </div>
            </section>
        </>
    )
}

export default HomePage
