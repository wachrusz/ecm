class FooterComponent {
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
        const footer = document.getElementById('footer');
        if (!footer) return;
        
        footer.innerHTML = `
            <div class="container">
                <div class="footer-content">
                    <div class="footer-section">
                        <h3>${this.state.locale === 'ru' ? 'Компания' : 
                            this.state.locale === 'en' ? 'Company' : 'Firma'}</h3>
                        <ul>
                            <li><a href="/about" onclick="router.navigate('/about'); return false;">
                                ${this.state.locale === 'ru' ? 'О компании' : 
                                  this.state.locale === 'en' ? 'About us' : 'O nas'}
                            </a></li>
                            <li><a href="/certificates" onclick="router.navigate('/certificates'); return false;">
                                ${this.state.locale === 'ru' ? 'Сертификаты' : 
                                  this.state.locale === 'en' ? 'Certificates' : 'Certyfikaty'}
                            </a></li>
                            <li><a href="/privacy" onclick="router.navigate('/privacy'); return false;">
                                ${this.state.locale === 'ru' ? 'Политика конфиденциальности' : 
                                  this.state.locale === 'en' ? 'Privacy policy' : 'Polityka prywatności'}
                            </a></li>
                        </ul>
                    </div>
                    
                    <div class="footer-section">
                        <h3>${this.state.locale === 'ru' ? 'Продукция' : 
                            this.state.locale === 'en' ? 'Products' : 'Produkty'}</h3>
                        <ul>
                            <li><a href="/products" onclick="router.navigate('/products'); return false;">
                                ${this.state.locale === 'ru' ? 'Все продукты' : 
                                  this.state.locale === 'en' ? 'All products' : 'Wszystkie produkty'}
                            </a></li>
                            <li><a href="/products/main-equipment" onclick="router.navigate('/products/main-equipment'); return false;">
                                ${this.state.locale === 'ru' ? 'Основное оборудование' : 
                                  this.state.locale === 'en' ? 'Main equipment' : 'Główne wyposażenie'}
                            </a></li>
                            <li><a href="/products/components" onclick="router.navigate('/products/components'); return false;">
                                ${this.state.locale === 'ru' ? 'Комплектующие' : 
                                  this.state.locale === 'en' ? 'Components' : 'Komponenty'}
                            </a></li>
                        </ul>
                    </div>
                    
                    <div class="footer-section">
                        <h3>${this.state.locale === 'ru' ? 'Контакты' : 
                            this.state.locale === 'en' ? 'Contacts' : 'Kontakty'}</h3>
                        <ul>
                            <li><i class="fas fa-map-marker-alt"></i> г. Москва, ул. Промышленная, д. 15</li>
                            <li><i class="fas fa-phone"></i> +7 (495) 123-45-67</li>
                            <li><i class="fas fa-envelope"></i> info@company.com</li>
                        </ul>
                    </div>
                    
                    <div class="footer-section">
                        <h3>${this.state.locale === 'ru' ? 'Подписка' : 
                            this.state.locale === 'en' ? 'Subscribe' : 'Subskrypcja'}</h3>
                        <p>${this.state.locale === 'ru' ? 'Подпишитесь на новости' : 
                            this.state.locale === 'en' ? 'Subscribe to news' : 'Subskrybuj aktualności'}</p>
                        <div class="subscribe-form">
                            <input type="email" placeholder="Email" id="subscribeEmail">
                            <button class="btn btn-accent" id="subscribeBtn">
                                ${this.state.locale === 'ru' ? 'Подписаться' : 
                                  this.state.locale === 'en' ? 'Subscribe' : 'Subskrybuj'}
                            </button>
                        </div>
                    </div>
                </div>
                
                <div class="footer-bottom">
                    <p>&copy; ${new Date().getFullYear()} Промышленное оборудование. 
                    ${this.state.locale === 'ru' ? 'Все права защищены' : 
                      this.state.locale === 'en' ? 'All rights reserved' : 'Wszelkie prawa zastrzeżone'}</p>
                </div>
            </div>
        `;
        
        this.attachEvents();
    }
    
    attachEvents() {
        const subscribeBtn = document.getElementById('subscribeBtn');
        const subscribeEmail = document.getElementById('subscribeEmail');
        
        if (subscribeBtn && subscribeEmail) {
            subscribeBtn.addEventListener('click', () => {
                const email = subscribeEmail.value.trim();
                if (this.validateEmail(email)) {
                    alert(this.state.locale === 'ru' ? 'Спасибо за подписку!' :
                          this.state.locale === 'en' ? 'Thank you for subscribing!' :
                          'Dziękujemy za subskrypcję!');
                    subscribeEmail.value = '';
                } else {
                    alert(this.state.locale === 'ru' ? 'Введите корректный email' :
                          this.state.locale === 'en' ? 'Please enter a valid email' :
                          'Proszę podać prawidłowy email');
                }
            });
        }
    }
    
    validateEmail(email) {
        const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return re.test(email);
    }
}

if (typeof window !== 'undefined') {
    window.FooterComponent = FooterComponent;
    console.log('✅ FooterComponent exported to window');
}