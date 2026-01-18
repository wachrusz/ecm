class HeaderComponent {
    constructor(store) {
        this.store = store;
        this.state = {};
        this.init();
    }
    
    init() {
        this.store.subscribe((state) => {
            this.state = state;
            this.render();
        });
    }
    
    render() {
        const header = document.getElementById('header');
        if (!header) return;
        
        header.innerHTML = `
            <div class="header-top">
                <div class="container">
                    <div class="header-top-content">
                        <div class="contacts-top">
                            <span><i class="fas fa-phone"></i> +7 (495) 123-45-67</span>
                            <span><i class="fas fa-envelope"></i> info@company.com</span>
                        </div>
                        <div class="language-switcher">
                            <button class="lang-btn ${this.state.locale === 'ru' ? 'active' : ''}" data-lang="ru">RU</button>
                            <button class="lang-btn ${this.state.locale === 'en' ? 'active' : ''}" data-lang="en">EN</button>
                            <button class="lang-btn ${this.state.locale === 'pl' ? 'active' : ''}" data-lang="pl">PL</button>
                        </div>
                    </div>
                </div>
            </div>
            <div class="header-main">
                <div class="container">
                    <div class="header-content">
                        <a href="/" class="logo" onclick="router.navigate('/'); return false;">
                            <span>Пром</span>Оборудование
                        </a>
                        
                        <button class="mobile-menu-btn">
                            <i class="fas fa-bars"></i>
                        </button>
                        
                        <nav>
                            <ul>
                                <li><a href="/" onclick="router.navigate('/'); return false;">
                                    ${this.state.locale === 'ru' ? 'Главная' : 
                                      this.state.locale === 'en' ? 'Home' : 'Strona główna'}
                                </a></li>
                                <li><a href="/about" onclick="router.navigate('/about'); return false;">
                                    ${this.state.locale === 'ru' ? 'О компании' : 
                                      this.state.locale === 'en' ? 'About' : 'O nas'}
                                </a></li>
                                <li><a href="/products" onclick="router.navigate('/products'); return false;">
                                    ${this.state.locale === 'ru' ? 'Продукция' : 
                                      this.state.locale === 'en' ? 'Products' : 'Produkty'}
                                </a></li>
                                <li><a href="/news" onclick="router.navigate('/news'); return false;">
                                    ${this.state.locale === 'ru' ? 'Новости' : 
                                      this.state.locale === 'en' ? 'News' : 'Aktualności'}
                                </a></li>
                                <li><a href="/documents" onclick="router.navigate('/documents'); return false;">
                                    ${this.state.locale === 'ru' ? 'Документы' : 
                                      this.state.locale === 'en' ? 'Documents' : 'Dokumenty'}
                                </a></li>
                                <li><a href="/contacts" onclick="router.navigate('/contacts'); return false;">
                                    ${this.state.locale === 'ru' ? 'Контакты' : 
                                      this.state.locale === 'en' ? 'Contacts' : 'Kontakty'}
                                </a></li>
                            </ul>
                        </nav>
                        
                        <div class="search-box">
                            <input type="text" id="searchInput" placeholder="${this.state.locale === 'ru' ? 'Поиск...' : 
                                this.state.locale === 'en' ? 'Search...' : 'Szukaj...'}">
                            <button id="searchBtn"><i class="fas fa-search"></i></button>
                        </div>
                    </div>
                </div>
            </div>
        `;
        
        this.attachEvents();
    }
    
    attachEvents() {
        document.querySelectorAll('.lang-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const lang = e.target.dataset.lang;
                this.store.changeLocale(lang);
                router.navigate(window.location.pathname);
            });
        });
        
        const mobileBtn = document.querySelector('.mobile-menu-btn');
        const nav = document.querySelector('nav');
        if (mobileBtn && nav) {
            mobileBtn.addEventListener('click', () => {
                nav.classList.toggle('active');
            });
        }
        
        // Поиск
        const searchBtn = document.getElementById('searchBtn');
        const searchInput = document.getElementById('searchInput');
        if (searchBtn && searchInput) {
            searchBtn.addEventListener('click', () => {
                const query = searchInput.value.trim();
                if (query) {
                    router.navigate(`/search?q=${encodeURIComponent(query)}`);
                }
            });
            
            searchInput.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    const query = searchInput.value.trim();
                    if (query) {
                        router.navigate(`/search?q=${encodeURIComponent(query)}`);
                    }
                }
            });
        }
    }
}