if (!window.Models || !window.HeaderComponent || !window.FooterComponent || !window.Router) {
    console.error('❌ Missing required classes. Check script loading order.');
    console.error('📊 Current status:');
    console.error('  - window.Models:', window.Models ? '✅ LOADED' : '❌ MISSING');
    console.error('  - window.HeaderComponent:', window.HeaderComponent ? '✅ LOADED' : '❌ MISSING');
    console.error('  - window.FooterComponent:', window.FooterComponent ? '✅ LOADED' : '❌ MISSING');
    console.error('  - window.Router:', window.Router ? '✅ LOADED' : '❌ MISSING');
    
    console.error('🔍 Checking script paths:');
    const scripts = document.querySelectorAll('script[src]');
    scripts.forEach(script => {
        console.error(`  - ${script.src}`);
    });
    
    document.addEventListener('DOMContentLoaded', function() {
        const content = document.getElementById('page-content');
        if (content) {
            content.innerHTML = `
                <div class="error" style="padding: 20px; background: #ffebee; border-radius: 8px; margin: 20px;">
                    <h2 style="color: #c62828;">⚠️ Ошибка загрузки приложения</h2>
                    <p>Отсутствуют необходимые компоненты:</p>
                    <ul style="text-align: left; display: inline-block;">
                        <li>Models: ${window.Models ? '✅' : '❌'}</li>
                        <li>HeaderComponent: ${window.HeaderComponent ? '✅' : '❌'}</li>
                        <li>FooterComponent: ${window.FooterComponent ? '✅' : '❌'}</li>
                        <li>Router: ${window.Router ? '✅' : '❌'}</li>
                    </ul>
                    <p>Проверьте консоль браузера (F12 → Console) для деталей</p>
                    <button onclick="location.reload()" 
                            style="padding: 10px 20px; background: #2196f3; color: white; border: none; border-radius: 4px; cursor: pointer; margin: 10px;">
                        🔄 Обновить страницу
                    </button>
                    <button onclick="checkScripts()"
                            style="padding: 10px 20px; background: #ff9800; color: white; border: none; border-radius: 4px; cursor: pointer; margin: 10px;">
                        🔍 Проверить скрипты
                    </button>
                </div>
            `;
        }
    });
} else {
    document.addEventListener('DOMContentLoaded', function() {
        try {
            console.log('Initializing app...');
            
            const store = new Models.Store();
            
            const savedLocale = localStorage.getItem('locale');
            if (savedLocale) {
                store.changeLocale(savedLocale);
            }
            
            new HeaderComponent(store);
            new FooterComponent(store);
            
            window.router = new Router(store);
            window.appStore = store;
            
            window.openModal = function(imageUrl) {
                const modal = document.getElementById('imageModal');
                const modalImg = document.getElementById('modalImage');
                if (modal && modalImg) {
                    modal.style.display = 'block';
                    modalImg.src = imageUrl;
                }
            };
            
            setTimeout(() => {
                store.setState({ isLoading: false });
                console.log('App initialized successfully');
            }, 500);
            
        } catch (error) {
            console.error('Error initializing app:', error);
            
            const content = document.getElementById('page-content');
            if (content) {
                content.innerHTML = `
                    <div class="error">
                        <h2>Ошибка инициализации приложения</h2>
                        <p>${error.message}</p>
                        <button onclick="location.reload()">Обновить страницу</button>
                    </div>
                `;
            }
        }
    });
}