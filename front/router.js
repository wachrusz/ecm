class Router {
    constructor(store) {
        this.store = store;
        this.routes = {
            '/': this.loadHomePage.bind(this),
            '/about': this.loadAboutPage.bind(this),
            '/certificates': this.loadCertificatesPage.bind(this),
            '/products': this.loadProductsPage.bind(this),
            '/products/:category': this.loadProductsByCategory.bind(this),
            '/product/:id': this.loadProductDetail.bind(this),
            '/news': this.loadNewsPage.bind(this),
            '/news/:id': this.loadNewsDetail.bind(this),
            '/documents': this.loadDocumentsPage.bind(this),
            '/contacts': this.loadContactsPage.bind(this),
            '/search': this.loadSearchPage.bind(this),
            '/privacy': this.loadPrivacyPage.bind(this)
        };
        
        this.init();
    }
    
    init() {
        document.addEventListener('click', (e) => {
            if (e.target.tagName === 'A' && e.target.href) {
                const url = new URL(e.target.href);
                if (url.origin === window.location.origin) {
                    e.preventDefault();
                    this.navigate(url.pathname + url.search);
                }
            }
        });
        
        window.addEventListener('popstate', () => {
            this.handleRoute();
        });
        
        this.handleRoute();
    }
    
    navigate(path) {
        window.history.pushState({}, '', path);
        this.handleRoute();
    }
    
    handleRoute() {
        const path = window.location.pathname;
        const search = window.location.search;
        const url = path + search;
        
        const localeMatch = path.match(/^\/(ru|en|pl)/);
        const locale = localeMatch ? localeMatch[1] : this.store.state.locale;
        
        this.store.changeLocale(locale);
        
        const routePath = localeMatch ? path.replace(`/${locale}`, '') || '/' : path;

        let matchedRoute = null;
        let params = {};

        for (const route in this.routes) {
            const routePattern = route.replace(/:\w+/g, '([^/]+)');
            const regex = new RegExp(`^${routePattern}$`);
            const match = routePath.match(regex);
            
            if (match) {
                matchedRoute = route;
                const paramNames = [...route.matchAll(/:(\w+)/g)].map(m => m[1]);
                paramNames.forEach((name, index) => {
                    params[name] = match[index + 1];
                });
                break;
            }
        }
        
        if (matchedRoute && this.routes[matchedRoute]) {
            this.store.setState({ isLoading: true });
            this.routes[matchedRoute](params, search);
        } else {
            this.load404Page();
        }
    }
    
    async loadHomePage() {
        await this.loadPage('home');
    }
    
    async loadAboutPage() {
        await this.loadPage('about');
    }
    
    async loadCertificatesPage() {
        await this.loadPage('certificates');
    }
    
    async loadProductsPage() {
        this.renderTemplate('products', async () => {
            try {
                const locale = this.store.state.locale;
                const response = await fetch(`/api/${locale}/products`);
                const data = await response.json();
                
                console.log('📦 Products API response:', data);
                
                const productsData = data.items || [];
                const products = productsData.map(item => new Models.Product(item));
                
                const categoriesMap = {};
                productsData.forEach(product => {
                    if (product.category && product.category.id) {
                        categoriesMap[product.category.id] = product.category;
                    }
                });
                
                const categories = Object.values(categoriesMap).map(cat => new Models.ProductCategory(cat));
                
                console.log('🎯 Products:', products);
                console.log('📂 Categories:', categories);
                
                this.store.setState({ products, categories });
                
                return `
                    <h1>${locale === 'ru' ? 'Продукция' : 
                        locale === 'en' ? 'Products' : 'Produkty'}</h1>
                    
                    <div class="categories-filter">
                        <button class="btn ${!this.store.state.selectedCategory ? 'active' : ''}" 
                                onclick="router.navigate('/products')">
                            ${locale === 'ru' ? 'Все' : 
                            locale === 'en' ? 'All' : 'Wszystkie'}
                        </button>
                        ${categories.map(cat => `
                            <button class="btn ${this.store.state.selectedCategory === cat.id ? 'active' : ''}" 
                                    onclick="router.navigate('/products/${cat.id}')">
                                ${cat.getName(locale)}
                            </button>
                        `).join('')}
                    </div>
                    
                    <div class="products-grid">
                        ${products.map(product => `
                            <div class="product-card card" onclick="router.navigate('/product/${product.id}')">
                                ${product.imageUrl ? `
                                    <div class="product-image">
                                        <img src="${product.imageUrl}" alt="${product.getName(locale)}">
                                    </div>
                                ` : ''}
                                <div class="product-info">
                                    <h3 class="product-title">${product.getName(locale)}</h3>
                                    <div class="product-sku">${product.sku}</div>
                                    ${product.getDescription(locale) ? `
                                        <p class="product-description">${product.getDescription(locale).substring(0, 100)}...</p>
                                    ` : ''}
                                    <button class="btn btn-primary">
                                        ${locale === 'ru' ? 'Подробнее' : 
                                        locale === 'en' ? 'Details' : 'Szczegóły'}
                                    </button>
                                </div>
                            </div>
                        `).join('')}
                    </div>
                    
                    ${products.length === 0 ? `
                        <div class="no-results">
                            <p>${locale === 'ru' ? 'Продукты не найдены' : 
                                locale === 'en' ? 'No products found' : 'Nie znaleziono produktów'}</p>
                        </div>
                    ` : ''}
                `;
            } catch (error) {
                console.error('Error loading products:', error);
                return `
                    <div class="error">
                        <h2>${this.store.state.locale === 'ru' ? 'Ошибка загрузки' : 
                            this.store.state.locale === 'en' ? 'Loading error' : 'Błąd ładowania'}</h2>
                        <p>${error.message}</p>
                    </div>
                `;
            }
        });
    }
    async loadProductsByCategory(params) {
        const { category } = params;
        this.store.setState({ selectedCategory: category });
        await this.loadProductsPage();
    }
    
    async loadProductDetail(params) {
        const { id } = params;
        
        this.renderTemplate('product-detail', async () => {
            try {
                const locale = this.store.state.locale;
                const response = await fetch(`/api/${locale}/product/${id}`);
                const data = await response.json();
                
                const product = new Models.Product(data);
                
                return `
                    <div class="breadcrumbs">
                        <a href="/" onclick="router.navigate('/'); return false;">
                            ${locale === 'ru' ? 'Главная' : 
                             locale === 'en' ? 'Home' : 'Strona główna'}
                        </a> / 
                        <a href="/products" onclick="router.navigate('/products'); return false;">
                            ${locale === 'ru' ? 'Продукция' : 
                             locale === 'en' ? 'Products' : 'Produkty'}
                        </a> / 
                        <span>${product.getName(locale)}</span>
                    </div>
                    
                    <div class="product-detail">
                        <div class="product-gallery">
                            ${product.imageUrl ? `
                                <div class="main-image">
                                    <img src="${product.imageUrl}" alt="${product.getName(locale)}" 
                                         onclick="openModal('${product.imageUrl}')">
                                </div>
                            ` : ''}
                        </div>
                        
                        <div class="product-info">
                            <h1>${product.getName(locale)}</h1>
                            <div class="product-sku">Артикул: ${product.sku}</div>
                            
                            ${product.getDescription(locale) ? `
                                <div class="product-description">
                                    <h3>${locale === 'ru' ? 'Описание' : 
                                         locale === 'en' ? 'Description' : 'Opis'}</h3>
                                    <p>${product.getDescription(locale)}</p>
                                </div>
                            ` : ''}
                            
                            ${product.fileUrl ? `
                                <a href="${product.fileUrl}" class="btn btn-secondary" target="_blank">
                                    <i class="fas fa-download"></i>
                                    ${locale === 'ru' ? 'Скачать техпаспорт' : 
                                     locale === 'en' ? 'Download datasheet' : 'Pobierz kartę katalogową'}
                                </a>
                            ` : ''}
                            
                            <button class="btn btn-accent contact-btn" 
                                    onclick="router.navigate('/contacts')">
                                <i class="fas fa-envelope"></i>
                                ${locale === 'ru' ? 'Запросить цену' : 
                                 locale === 'en' ? 'Request price' : 'Zapytaj o cenę'}
                            </button>
                        </div>
                    </div>
                    
                    ${product.specs && product.specs.length > 0 ? `
                        <div class="product-specs">
                            <h2>${locale === 'ru' ? 'Характеристики' : 
                                  locale === 'en' ? 'Specifications' : 'Specyfikacje'}</h2>
                            <table class="spec-table">
                                ${product.specs.map(spec => `
                                    <tr>
                                        <td>${spec.name}</td>
                                        <td>${spec.value}</td>
                                    </tr>
                                `).join('')}
                            </table>
                        </div>
                    ` : ''}
                `;
            } catch (error) {
                console.error('Error loading product:', error);
                return `
                    <div class="error">
                        <h2>${this.store.state.locale === 'ru' ? 'Продукт не найден' : 
                              this.store.state.locale === 'en' ? 'Product not found' : 'Produkt nie znaleziony'}</h2>
                        <a href="/products" class="btn btn-primary" onclick="router.navigate('/products'); return false;">
                            ${this.store.state.locale === 'ru' ? 'Вернуться к продукции' : 
                             this.store.state.locale === 'en' ? 'Back to products' : 'Powrót do produktów'}
                        </a>
                    </div>
                `;
            }
        });
    }
    
    async loadNewsPage() {
        this.renderTemplate('news', async () => {
            try {
                const locale = this.store.state.locale;
                const response = await fetch(`/api/${locale}/news`);
                const data = await response.json();
                
                const newsData = data.items || [];
                const newsList = newsData.map(item => new Models.News(item));
                
                return `
                    <h1>${locale === 'ru' ? 'Новости' : 
                          locale === 'en' ? 'News' : 'Aktualności'}</h1>
                    
                    <div class="news-grid">
                        ${newsList.map(news => `
                            <div class="news-card card" onclick="router.navigate('/news/${news.id}')">
                                ${news.imageUrl ? `
                                    <div class="news-image">
                                        <img src="${news.imageUrl}" alt="${news.getTitle(locale)}">
                                    </div>
                                ` : ''}
                                <div class="news-info">
                                    <div class="news-date">
                                        ${new Date(news.createdAt).toLocaleDateString(locale)}
                                    </div>
                                    <h3 class="news-title">${news.getTitle(locale)}</h3>
                                    <p class="news-excerpt">${news.getContent(locale).substring(0, 150)}...</p>
                                    <button class="btn btn-primary">
                                        ${locale === 'ru' ? 'Читать далее' : 
                                         locale === 'en' ? 'Read more' : 'Czytaj więcej'}
                                    </button>
                                </div>
                            </div>
                        `).join('')}
                    </div>
                `;
            } catch (error) {
                console.error('Error loading news:', error);
                return `
                    <div class="error">
                        <h2>${this.store.state.locale === 'ru' ? 'Ошибка загрузки новостей' : 
                              this.store.state.locale === 'en' ? 'Error loading news' : 'Błąd ładowania aktualności'}</h2>
                    </div>
                `;
            }
        });
    }
    
    async loadNewsDetail(params) {
        const { id } = params;
        
        this.renderTemplate('news-detail', async () => {
            try {
                const locale = this.store.state.locale;
                this.store.state.News
                const response = await fetch(`/api/${locale}/news/${id}`);
                const data = await response.json();
                
                const news = new Models.News(data);
                
                return `
                    <div class="breadcrumbs">
                        <a href="/" onclick="router.navigate('/'); return false;">
                            ${locale === 'ru' ? 'Главная' : 
                             locale === 'en' ? 'Home' : 'Strona główna'}
                        </a> / 
                        <a href="/news" onclick="router.navigate('/news'); return false;">
                            ${locale === 'ru' ? 'Новости' : 
                             locale === 'en' ? 'News' : 'Aktualności'}
                        </a> / 
                        <span>${news.getTitle(locale)}</span>
                    </div>
                    
                    <article class="news-article">
                        <div class="news-header">
                            <h1>${news.getTitle(locale)}</h1>
                            <div class="news-meta">
                                <span class="news-date">
                                    <i class="far fa-calendar"></i>
                                    ${new Date(news.createdAt).toLocaleDateString(locale, {
                                        year: 'numeric',
                                        month: 'long',
                                        day: 'numeric'
                                    })}
                                </span>
                            </div>
                        </div>
                        
                        ${news.imageUrl ? `
                            <div class="news-image-full">
                                <img src="${news.imageUrl}" alt="${news.getTitle(locale)}">
                            </div>
                        ` : ''}
                        
                        <div class="news-content">
                            ${news.getContent(locale)}
                        </div>
                        
                        <div class="news-footer">
                            <a href="/news" class="btn btn-primary" onclick="router.navigate('/news'); return false;">
                                <i class="fas fa-arrow-left"></i>
                                ${locale === 'ru' ? 'Все новости' : 
                                 locale === 'en' ? 'All news' : 'Wszystkie aktualności'}
                            </a>
                        </div>
                    </article>
                `;
            } catch (error) {
                console.error('Error loading news detail:', error);
                return `
                    <div class="error">
                        <h2>${this.store.state.locale === 'ru' ? 'Новость не найдена' : 
                              this.store.state.locale === 'en' ? 'News not found' : 'Aktualność nie znaleziona'}</h2>
                    </div>
                `;
            }
        });
    }
    
    async loadDocumentsPage() {
        this.renderTemplate('documents', async () => {
            try {
                const locale = this.store.state.locale;
                const response = await fetch(`/api/${locale}/documents`);
                const data = await response.json();
                
                const documents = data.map(item => new Models.Document(item));
                
                const grouped = {};
                documents.forEach(doc => {
                    if (!grouped[doc.type]) grouped[doc.type] = [];
                    grouped[doc.type].push(doc);
                });
                
                const typeLabels = {
                    'ru': { gost: 'ГОСТы', certificate: 'Сертификаты', reference: 'Справочники' },
                    'en': { gost: 'GOST Standards', certificate: 'Certificates', reference: 'References' },
                    'pl': { gost: 'Normy GOST', certificate: 'Certyfikaty', reference: 'Referencje' }
                };
                
                return `
                    <h1>${locale === 'ru' ? 'Документация' : 
                          locale === 'en' ? 'Documents' : 'Dokumentacja'}</h1>
                    
                    ${Object.keys(grouped).map(type => `
                        <div class="document-section">
                            <h2>${typeLabels[locale][type] || type}</h2>
                            <div class="documents-list">
                                ${grouped[type].map(doc => `
                                    <div class="document-item card">
                                        <div class="document-info">
                                            <h3>${doc.getTitle(locale)}</h3>
                                            <div class="document-date">
                                                ${new Date(doc.createdAt).toLocaleDateString(locale)}
                                            </div>
                                        </div>
                                        <a href="${doc.fileUrl}" class="btn btn-secondary" target="_blank">
                                            <i class="fas fa-download"></i>
                                            ${locale === 'ru' ? 'Скачать' : 
                                             locale === 'en' ? 'Download' : 'Pobierz'}
                                        </a>
                                    </div>
                                `).join('')}
                            </div>
                        </div>
                    `).join('')}
                `;
            } catch (error) {
                console.error('Error loading documents:', error);
                return `
                    <div class="error">
                        <h2>${this.store.state.locale === 'ru' ? 'Ошибка загрузки документов' : 
                              this.store.state.locale === 'en' ? 'Error loading documents' : 'Błąd ładowania dokumentów'}</h2>
                    </div>
                `;
            }
        });
    }
    
    async loadContactsPage() {
        this.renderTemplate('contacts', async () => {
            try {
                const locale = this.store.state.locale;
                const response = await fetch(`/api/${locale}/contacts`);
                const data = await response.json();
                
                const contacts = data.contacts?.map(item => new Models.Contact(item)) || [];
                
                return `
                    <h1>${locale === 'ru' ? 'Контакты' : 
                          locale === 'en' ? 'Contacts' : 'Kontakty'}</h1>
                    
                    <div class="contacts-page">
                        <div class="contact-info">
                            ${contacts.map(contact => `
                                <div class="contact-item">
                                    <i class="fas fa-${this.getContactIcon(contact.type)}"></i>
                                    <div>
                                        <strong>${contact.getLabel(locale)}:</strong>
                                        <p>${contact.value}</p>
                                    </div>
                                </div>
                            `).join('')}
                            
                            <div class="map-container">
                                <iframe src="https://yandex.ru/map-widget/v1/?um=constructor%3A12345&amp;source=constructor" 
                                        width="100%" height="400" frameborder="0"></iframe>
                            </div>
                        </div>
                        
                        <div class="contact-form-container">
                            <h2>${locale === 'ru' ? 'Написать сообщение' : 
                                  locale === 'en' ? 'Send message' : 'Wyślij wiadomość'}</h2>
                            <form id="feedbackForm">
                                <div class="form-group">
                                    <label for="name">${locale === 'ru' ? 'Имя' : 
                                                      locale === 'en' ? 'Name' : 'Imię'} *</label>
                                    <input type="text" id="name" name="name" class="form-control" required>
                                </div>
                                
                                <div class="form-group">
                                    <label for="email">Email *</label>
                                    <input type="email" id="email" name="email" class="form-control" required>
                                </div>
                                
                                <div class="form-group">
                                    <label for="phone">${locale === 'ru' ? 'Телефон' : 
                                                       locale === 'en' ? 'Phone' : 'Telefon'}</label>
                                    <input type="tel" id="phone" name="phone" class="form-control">
                                </div>
                                
                                <div class="form-group">
                                    <label for="company">${locale === 'ru' ? 'Компания' : 
                                                         locale === 'en' ? 'Company' : 'Firma'}</label>
                                    <input type="text" id="company" name="company" class="form-control">
                                </div>
                                
                                <div class="form-group">
                                    <label for="message">${locale === 'ru' ? 'Сообщение' : 
                                                         locale === 'en' ? 'Message' : 'Wiadomość'} *</label>
                                    <textarea id="message" name="message" class="form-control" required></textarea>
                                </div>
                                
                                <button type="submit" class="btn btn-primary">
                                    ${locale === 'ru' ? 'Отправить' : 
                                     locale === 'en' ? 'Send' : 'Wyślij'}
                                </button>
                            </form>
                        </div>
                    </div>
                `;
            } catch (error) {
                console.error('Error loading contacts:', error);
                return `
                    <div class="error">
                        <h2>${this.store.state.locale === 'ru' ? 'Ошибка загрузки контактов' : 
                              this.store.state.locale === 'en' ? 'Error loading contacts' : 'Błąd ładowania kontaktów'}</h2>
                    </div>
                `;
            } finally {
                setTimeout(() => this.initFeedbackForm(), 100);
            }
        });
    }
    
    async loadSearchPage() {
        const urlParams = new URLSearchParams(window.location.search);
        const query = urlParams.get('q') || '';
        
        this.renderTemplate('search', async () => {
            try {
                const locale = this.store.state.locale;
                const response = await fetch(`/api/${locale}/search?q=${encodeURIComponent(query)}`);
                const data = await response.json();
                
                return `
                    <h1>${locale === 'ru' ? 'Результаты поиска' : 
                          locale === 'en' ? 'Search results' : 'Wyniki wyszukiwania'}</h1>
                    
                    <div class="search-info">
                        <p>${locale === 'ru' ? 'Найдено' : 
                            locale === 'en' ? 'Found' : 'Znaleziono'} ${data.count || 0} 
                            ${locale === 'ru' ? 'результатов по запросу' : 
                             locale === 'en' ? 'results for' : 'wyników dla'} "${query}"</p>
                    </div>
                    
                    ${data.results && data.results.length > 0 ? `
                        <div class="search-results">
                            ${data.results.map(result => `
                                <div class="search-result card">
                                    <h3><a href="${result.url}">${result.title}</a></h3>
                                    <p>${result.excerpt}</p>
                                    <div class="result-type">${result.type}</div>
                                </div>
                            `).join('')}
                        </div>
                    ` : `
                        <div class="no-results">
                            <p>${locale === 'ru' ? 'По вашему запросу ничего не найдено' : 
                                locale === 'en' ? 'No results found for your query' : 
                                'Nie znaleziono wyników dla Twojego zapytania'}</p>
                        </div>
                    `}
                `;
            } catch (error) {
                console.error('Error loading search results:', error);
                return `
                    <div class="error">
                        <h2>${this.store.state.locale === 'ru' ? 'Ошибка поиска' : 
                              this.store.state.locale === 'en' ? 'Search error' : 'Błąd wyszukiwania'}</h2>
                    </div>
                `;
            }
        });
    }
    
    async loadPrivacyPage() {
        await this.loadPage('privacy');
    }
    
    async loadPage(slug) {
        this.renderTemplate('page', async () => {
            try {
                const locale = this.store.state.locale;
                const response = await fetch(`/api/${locale}/page/${slug}`);
                const data = await response.json();
                console.log(data)
                
                const page = new Models.Page(data);
                const translation = page.getTranslation(locale);
                
                let contentHtml = '';
                if (translation.content) {
                    try {
                        const content = JSON.parse(translation.content);
                        console.log("content", content)

                        contentHtml = this.renderPageContent(content, locale);
                    } catch (e) {
                        contentHtml = `<div class="page-content">${translation.content}</div>`;
                    }
                }
                
                return `
                    <h1>${translation.title || ''}</h1>
                    ${contentHtml}
                `;
            } catch (error) {
                console.error(`Error loading page ${slug}:`, error);
                return `
                    <div class="error">
                        <h2>${this.store.state.locale === 'ru' ? 'Страница не найдена' : 
                              this.store.state.locale === 'en' ? 'Page not found' : 'Strona nie znaleziona'}</h2>
                    </div>
                `;
            }
        });
    }
    
    load404Page() {
        const locale = this.store.state.locale;
        const content = document.getElementById('page-content');
        if (content) {
            content.innerHTML = `
                <div class="error-page">
                    <h1>404</h1>
                    <h2>${locale === 'ru' ? 'Страница не найдена' : 
                          locale === 'en' ? 'Page not found' : 'Strona nie znaleziona'}</h2>
                    <p>${locale === 'ru' ? 'Запрашиваемая страница не существует' : 
                        locale === 'en' ? 'The requested page does not exist' : 
                        'Żądana strona nie istnieje'}</p>
                    <a href="/" class="btn btn-primary" onclick="router.navigate('/'); return false;">
                        ${locale === 'ru' ? 'На главную' : 
                         locale === 'en' ? 'Go to homepage' : 'Przejdź do strony głównej'}
                    </a>
                </div>
            `;
            this.store.setState({ isLoading: false });
        }
    }
    
    renderPageContent(content, locale) {
        console.warn("Calling render page content for: ", content, locale)

        let html = '';
        
        if (content.hero) {
            html += `
                <div class="hero-section">
                    <h2>${content.hero.title || ''}</h2>
                    <p>${content.hero.subtitle || ''}</p>
                    ${content.hero.cta ? `
                        <a href="/products" class="btn btn-accent" onclick="router.navigate('/products'); return false;">
                            ${content.hero.cta}
                        </a>
                    ` : ''}
                </div>
            `;
        }
        
        if (content.features && Array.isArray(content.features)) {
            html += `
                <div class="features-section">
                    <h3>${locale === 'ru' ? 'Наши преимущества' : 
                          locale === 'en' ? 'Our advantages' : 'Nasze zalety'}</h3>
                    <div class="features-grid">
                        ${content.features.map(feature => `
                            <div class="feature-card">
                                <h4>${feature.title}</h4>
                                <p>${feature.description}</p>
                            </div>
                        `).join('')}
                    </div>
                </div>
            `;
        }
        
        if (content.about) {
            html += `
                <div class="about-section">
                    <h2>${content.about.title || ''}</h2>
                    <p>${content.about.content || ''}</p>
                </div>
            `;
        }
        
        if (content.stats && Array.isArray(content.stats)) {
            html += `
                <div class="stats-section">
                    ${content.stats.map(stat => `
                        <div class="stat-item">
                            <div class="stat-value">${stat.value}</div>
                            <div class="stat-label">${stat.label}</div>
                        </div>
                    `).join('')}
                </div>
            `;
        }
        
        return html;
    }
    
    renderTemplate(templateName, contentCallback) {
        const content = document.getElementById('page-content');
        if (content) {
            content.innerHTML = `
                <div class="loading">
                    <i class="fas fa-spinner fa-spin"></i>
                    <p>${this.store.state.locale === 'ru' ? 'Загрузка...' : 
                         this.store.state.locale === 'en' ? 'Loading...' : 'Ładowanie...'}</p>
                </div>
            `;
            
            contentCallback().then(html => {
                content.innerHTML = `<div class="fade-in">${html}</div>`;
                this.store.setState({ isLoading: false });
            }).catch(error => {
                content.innerHTML = `
                    <div class="error">
                        <h2>${this.store.state.locale === 'ru' ? 'Ошибка загрузки' : 
                              this.store.state.locale === 'en' ? 'Loading error' : 'Błąd ładowania'}</h2>
                        <p>${error.message}</p>
                    </div>
                `;
                this.store.setState({ isLoading: false });
            });
        }
    }
    
    getContactIcon(type) {
        const icons = {
            'phone': 'phone',
            'email': 'envelope',
            'address': 'map-marker-alt',
            'map': 'map'
        };
        return icons[type] || 'info-circle';
    }
    
    initFeedbackForm() {
        const form = document.getElementById('feedbackForm');
        if (form) {
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                
                const formData = new FormData(form);
                const data = Object.fromEntries(formData.entries());
                
                try {
                    const locale = this.store.state.locale;
                    const response = await fetch(`/api/${locale}/feedback`, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify(data)
                    });
                    
                    if (response.ok) {
                        alert(locale === 'ru' ? 'Сообщение отправлено!' :
                              locale === 'en' ? 'Message sent!' :
                              'Wiadomość wysłana!');
                        form.reset();
                    } else {
                        throw new Error('Failed to send message');
                    }
                } catch (error) {
                    console.error('Error sending feedback:', error);
                    alert(locale === 'ru' ? 'Ошибка при отправке сообщения' :
                          locale === 'en' ? 'Error sending message' :
                          'Błąd podczas wysyłania wiadomości');
                }
            });
        }
    }
}

window.openModal = function(imageUrl) {
    const modal = document.getElementById('imageModal');
    const modalImg = document.getElementById('modalImage');
    if (modal && modalImg) {
        modal.style.display = 'block';
        modalImg.src = imageUrl;
    }
};

document.addEventListener('DOMContentLoaded', function() {
    const modal = document.getElementById('imageModal');
    const closeBtn = document.querySelector('.modal .close');
    
    if (closeBtn) {
        closeBtn.addEventListener('click', () => {
            modal.style.display = 'none';
        });
    }
    
    window.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.style.display = 'none';
        }
    });
});

if (typeof window !== 'undefined') {
    window.Router = Router;
    console.log('✅ Router exported to window');
}