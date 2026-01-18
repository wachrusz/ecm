-- Enable UUID extension if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Languages
CREATE TABLE languages (
    code VARCHAR(10) PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- Insert default languages
INSERT INTO languages (code, name) VALUES
('ru', 'Русский'),
('en', 'English'),
('pl', 'Polski');

-- Pages
CREATE TABLE pages (
    page_id SERIAL PRIMARY KEY,
    slug VARCHAR(255) UNIQUE NOT NULL,
    template VARCHAR(100) NOT NULL DEFAULT 'default',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Page translations
CREATE TABLE page_translations (
    page_id INT REFERENCES pages(page_id) ON DELETE CASCADE,
    language_code VARCHAR(10) REFERENCES languages(code),
    title VARCHAR(255) NOT NULL,
    content TEXT,
    meta_title VARCHAR(255),
    meta_description TEXT,
    PRIMARY KEY (page_id, language_code)
);

-- Product categories
CREATE TABLE product_categories (
    category_id SERIAL PRIMARY KEY,
    parent_id INT REFERENCES product_categories(category_id),
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Product category translations
CREATE TABLE product_category_translations (
    category_id INT REFERENCES product_categories(category_id) ON DELETE CASCADE,
    language_code VARCHAR(10) REFERENCES languages(code),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    PRIMARY KEY (category_id, language_code)
);

-- Products
CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,
    category_id INT REFERENCES product_categories(category_id),
    sku VARCHAR(100) UNIQUE NOT NULL,
    image_url VARCHAR(500),
    file_url VARCHAR(500),
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Product translations
CREATE TABLE product_translations (
    product_id INT REFERENCES products(product_id) ON DELETE CASCADE,
    language_code VARCHAR(10) REFERENCES languages(code),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    short_description TEXT,
    PRIMARY KEY (product_id, language_code)
);

-- Product specifications
CREATE TABLE product_specs (
    spec_id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(product_id) ON DELETE CASCADE,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Product spec translations
CREATE TABLE product_spec_translations (
    spec_id INT REFERENCES product_specs(spec_id) ON DELETE CASCADE,
    language_code VARCHAR(10) REFERENCES languages(code),
    name VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY (spec_id, language_code)
);

-- News
CREATE TABLE news (
    news_id SERIAL PRIMARY KEY,
    image_url VARCHAR(500),
    published BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- News translations
CREATE TABLE news_translations (
    news_id INT REFERENCES news(news_id) ON DELETE CASCADE,
    language_code VARCHAR(10) REFERENCES languages(code),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT,
    PRIMARY KEY (news_id, language_code)
);

-- Documents
CREATE TABLE documents (
    document_id SERIAL PRIMARY KEY,
    file_url VARCHAR(500) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'gost', 'certificate', 'reference'
    created_at TIMESTAMP DEFAULT NOW()
);

-- Document translations
CREATE TABLE document_translations (
    document_id INT REFERENCES documents(document_id) ON DELETE CASCADE,
    language_code VARCHAR(10) REFERENCES languages(code),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    PRIMARY KEY (document_id, language_code)
);

-- Contacts
CREATE TABLE contacts (
    contact_id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL, -- 'phone', 'email', 'address', 'map'
    value TEXT NOT NULL,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Contact translations
CREATE TABLE contact_translations (
    contact_id INT REFERENCES contacts(contact_id) ON DELETE CASCADE,
    language_code VARCHAR(10) REFERENCES languages(code),
    label VARCHAR(255),
    PRIMARY KEY (contact_id, language_code)
);

-- Feedback
CREATE TABLE feedback (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    company VARCHAR(255),
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    processed BOOLEAN DEFAULT false
);

-- Indexes for better performance
CREATE INDEX idx_pages_slug ON pages(slug);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_news_published ON news(published, created_at);
CREATE INDEX idx_feedback_processed ON feedback(processed, created_at);

INSERT INTO pages (slug, template, created_at, updated_at) VALUES
('home', 'homepage', NOW(), NOW()),
('about', 'default', NOW(), NOW()),
('contacts', 'contacts', NOW(), NOW()),
('products', 'products', NOW(), NOW()),
('news', 'news', NOW(), NOW()),
('documents', 'documents', NOW(), NOW());

INSERT INTO page_translations (page_id, language_code, title, content, meta_title, meta_description) VALUES
(1, 'ru', 'Главная', '{"hero": {"title": "Производство промышленного оборудования", "subtitle": "Высококачественные решения для вашего бизнеса", "cta": "Посмотреть каталог"}, "features": [{"title": "Качество", "description": "Сертифицированная продукция"}, {"title": "Опыт", "description": "Более 20 лет на рынке"}, {"title": "Поддержка", "description": "Техническое сопровождение"}]}', 'Промышленное оборудование | Главная', 'Производство промышленного оборудования высокого качества. Сертифицированная продукция, техническая поддержка.'),
(1, 'en', 'Home', '{"hero": {"title": "Industrial Equipment Manufacturing", "subtitle": "High-quality solutions for your business", "cta": "View catalog"}, "features": [{"title": "Quality", "description": "Certified products"}, {"title": "Experience", "description": "Over 20 years on the market"}, {"title": "Support", "description": "Technical support"}]}', 'Industrial Equipment | Home', 'High-quality industrial equipment manufacturing. Certified products, technical support.'),
(1, 'pl', 'Strona główna', '{"hero": {"title": "Produkcja sprzętu przemysłowego", "subtitle": "Wysokiej jakości rozwiązania dla Twojego biznesu", "cta": "Zobacz katalog"}, "features": [{"title": "Jakość", "description": "Certyfikowane produkty"}, {"title": "Doświadczenie", "description": "Ponad 20 lat na rynku"}, {"title": "Wsparcie", "description": "Wsparcie techniczne"}]}', 'Sprzęt przemysłowy | Strona główna', 'Produkcja wysokiej jakości sprzętu przemysłowego. Produkty certyfikowane, wsparcie techniczne.'),
(2, 'ru', 'О компании', '{"about": {"title": "Наша история", "content": "Компания основана в 2000 году. Специализируемся на производстве промышленного оборудования."}, "stats": [{"value": "20+", "label": "Лет опыта"}, {"value": "500+", "label": "Проектов"}, {"value": "50+", "label": "Сотрудников"}]}', 'О компании | Производство оборудования', 'Информация о нашей компании, истории и достижениях.'),
(2, 'en', 'About Us', '{"about": {"title": "Our Story", "content": "Company founded in 2000. We specialize in industrial equipment manufacturing."}, "stats": [{"value": "20+", "label": "Years of experience"}, {"value": "500+", "label": "Projects"}, {"value": "50+", "label": "Employees"}]}', 'About Us | Equipment Manufacturing', 'Information about our company, history and achievements.'),
(2, 'pl', 'O nas', '{"about": {"title": "Nasza historia", "content": "Firma założona w 2000 roku. Specjalizujemy się w produkcji sprzętu przemysłowego."}, "stats": [{"value": "20+", "label": "Lat doświadczenia"}, {"value": "500+", "label": "Projektów"}, {"value": "50+", "label": "Pracowników"}]}', 'O nas | Produkcja sprzętu', 'Informacje o naszej firmie, historii i osiągnięciach.'),
(3, 'ru', 'Контакты', '{"contact_info": {"title": "Свяжитесь с нами", "description": "Мы всегда рады ответить на ваши вопросы"}, "form": {"title": "Написать сообщение", "fields": ["name", "email", "phone", "message"]}}', 'Контакты | Наши координаты', 'Контактная информация, адрес, телефоны, форма обратной связи.'),
(3, 'en', 'Contacts', '{"contact_info": {"title": "Contact us", "description": "We are always happy to answer your questions"}, "form": {"title": "Send message", "fields": ["name", "email", "phone", "message"]}}', 'Contacts | Our coordinates', 'Contact information, address, phones, feedback form.'),
(3, 'pl', 'Kontakty', '{"contact_info": {"title": "Skontaktuj się z nami", "description": "Zawsze chętnie odpowiemy na Twoje pytania"}, "form": {"title": "Wyślij wiadomość", "fields": ["name", "email", "phone", "message"]}}', 'Kontakty | Nasze dane', 'Dane kontaktowe, adres, telefony, formularz kontaktowy.'),
(4, 'ru', 'Продукция', '{"catalog": {"title": "Наша продукция", "filter": {"categories": true, "specifications": true}}, "layout": "grid"}', 'Каталог продукции | Промышленное оборудование', 'Каталог промышленного оборудования и комплектующих.'),
(4, 'en', 'Products', '{"catalog": {"title": "Our products", "filter": {"categories": true, "specifications": true}}, "layout": "grid"}', 'Product catalog | Industrial equipment', 'Catalog of industrial equipment and components.'),
(4, 'pl', 'Produkty', '{"catalog": {"title": "Nasze produkty", "filter": {"categories": true, "specifications": true}}, "layout": "grid"}', 'Katalog produktów | Sprzęt przemysłowy', 'Katalog sprzętu przemysłowego i komponentów.'),
(5, 'ru', 'Новости', '{"news": {"title": "Последние новости", "per_page": 10, "show_dates": true}}', 'Новости компании и отрасли', 'Актуальные новости о компании, продукции и промышленности.'),
(5, 'en', 'News', '{"news": {"title": "Latest news", "per_page": 10, "show_dates": true}}', 'Company and industry news', 'Current news about the company, products and industry.'),
(5, 'pl', 'Aktualności', '{"news": {"title": "Najnowsze aktualności", "per_page": 10, "show_dates": true}}', 'Aktualności firmy i branży', 'Aktualne informacje o firmie, produktach i branży.'),
(6, 'ru', 'Документы', '{"documents": {"title": "Техническая документация", "categories": ["ГОСТы", "Сертификаты", "Справочники"]}}', 'Техническая документация и сертификаты', 'ГОСТы, сертификаты качества, технические справочники.'),
(6, 'en', 'Documents', '{"documents": {"title": "Technical documentation", "categories": ["GOST", "Certificates", "References"]}}', 'Technical documentation and certificates', 'GOST standards, quality certificates, technical references.'),
(6, 'pl', 'Dokumenty', '{"documents": {"title": "Dokumentacja techniczna", "categories": ["GOST", "Certyfikaty", "Referencje"]}}', 'Dokumentacja techniczna i certyfikaty', 'Standardy GOST, certyfikaty jakości, referencje techniczne.');

INSERT INTO contacts (type, value, sort_order) VALUES
('phone', '+7 (495) 123-45-67', 1),
('phone', '+7 (800) 555-35-35', 2),
('email', 'info@company.com', 3),
('address', 'г. Москва, ул. Промышленная, д. 15', 4),
('map', 'https://yandex.ru/maps/-/CDVZIIB9', 5);

INSERT INTO contact_translations (contact_id, language_code, label) VALUES
(1, 'ru', 'Основной телефон'),
(1, 'en', 'Main phone'),
(1, 'pl', 'Telefon główny'),
(2, 'ru', 'Бесплатный номер'),
(2, 'en', 'Toll-free number'),
(2, 'pl', 'Numer bezpłatny'),
(3, 'ru', 'Электронная почта'),
(3, 'en', 'Email'),
(3, 'pl', 'E-mail'),
(4, 'ru', 'Адрес офиса'),
(4, 'en', 'Office address'),
(4, 'pl', 'Adres biura'),
(5, 'ru', 'Карта проезда'),
(5, 'en', 'Location map'),
(5, 'pl', 'Mapa dojazdu');

-- Добавление категорий продукции
INSERT INTO product_categories (parent_id, sort_order) VALUES
(NULL, 1), -- 1. Основное оборудование
(NULL, 2), -- 2. Комплектующие
(1, 1),    -- 3. Станки (дочерняя от 1)
(1, 2),    -- 4. Прессы (дочерняя от 1)
(2, 1);    -- 5. Электроника (дочерняя от 2)

INSERT INTO product_category_translations (category_id, language_code, name, description) VALUES
(1, 'ru', 'Основное оборудование', 'Основные производственные линии и станки'),
(1, 'en', 'Main equipment', 'Main production lines and machines'),
(1, 'pl', 'Główne wyposażenie', 'Główne linie produkcyjne i maszyny'),
(2, 'ru', 'Комплектующие', 'Запасные части и компоненты'),
(2, 'en', 'Components', 'Spare parts and components'),
(2, 'pl', 'Komponenty', 'Części zamienne i komponenty'),
(3, 'ru', 'Станки', 'Металлообрабатывающие станки'),
(3, 'en', 'Machine tools', 'Metalworking machines'),
(3, 'pl', 'Obrabiarki', 'Obrabiarki do metalu'),
(4, 'ru', 'Прессы', 'Гидравлические и механические прессы'),
(4, 'en', 'Presses', 'Hydraulic and mechanical presses'),
(4, 'pl', 'Prasy', 'Prasy hydrauliczne i mechaniczne'),
(5, 'ru', 'Электроника', 'Системы управления и автоматизации'),
(5, 'en', 'Electronics', 'Control and automation systems'),
(5, 'pl', 'Elektronika', 'Systemy sterowania i automatyki');

INSERT INTO products (category_id, sku, image_url, file_url, sort_order) VALUES
(3, 'CNC-1000', '/images/products/cnc-1000.jpg', '/files/manuals/cnc-1000.pdf', 1),
(3, 'MILL-500', '/images/products/mill-500.jpg', '/files/manuals/mill-500.pdf', 2),
(4, 'PRESS-H200', '/images/products/press-h200.jpg', '/files/manuals/press-h200.pdf', 1),
(5, 'CONTROL-X1', '/images/products/control-x1.jpg', '/files/manuals/control-x1.pdf', 1);

INSERT INTO product_translations (product_id, language_code, name, description, short_description) VALUES
(1, 'ru', 'Станок ЧПУ CNC-1000', 'Высокоточный станок с ЧПУ для металлообработки. Автоматическая смена инструмента, система охлаждения.', 'Станок ЧПУ для точной обработки'),
(1, 'en', 'CNC-1000 Machine', 'High-precision CNC machine for metalworking. Automatic tool changer, cooling system.', 'CNC machine for precision machining'),
(1, 'pl', 'Obrabiarka CNC-1000', 'Wysokiej precyzji obrabiarka CNC do obróbki metali. Automatyczna zmiana narzędzi, system chłodzenia.', 'Obrabiarka CNC do precyzyjnej obróbki'),
(2, 'ru', 'Фрезерный станок MILL-500', 'Универсальный фрезерный станок для обработки различных материалов.', 'Универсальный фрезерный станок'),
(2, 'en', 'MILL-500 Milling Machine', 'Universal milling machine for processing various materials.', 'Universal milling machine'),
(2, 'pl', 'Frezarka MILL-500', 'Uniwersalna frezarka do obróbki różnych materiałów.', 'Uniwersalna frezarka'),
(3, 'ru', 'Гидравлический пресс H-200', 'Мощный гидравлический пресс для штамповки и гибки металла.', 'Гидравлический пресс 200 тонн'),
(3, 'en', 'Hydraulic Press H-200', 'Powerful hydraulic press for metal stamping and bending.', 'Hydraulic press 200 tons'),
(3, 'pl', 'Prasa hydrauliczna H-200', 'Wydajna prasa hydrauliczna do tłoczenia i gięcia metalu.', 'Prasa hydrauliczna 200 ton'),
(4, 'ru', 'Система управления CONTROL-X1', 'Цифровая система управления для промышленного оборудования.', 'Система ЧПУ нового поколения'),
(4, 'en', 'CONTROL-X1 Control System', 'Digital control system for industrial equipment.', 'Next generation CNC system'),
(4, 'pl', 'System sterowania CONTROL-X1', 'Cyfrowy system sterowania dla urządzeń przemysłowych.', 'Nowa generacja systemu CNC');

INSERT INTO product_specs (product_id, sort_order) VALUES
(1, 1),
(1, 2),
(1, 3),
(1, 4);

INSERT INTO product_spec_translations (spec_id, language_code, name, value) VALUES
(1, 'ru', 'Макс. размер заготовки', '1000 × 800 × 600 мм'),
(1, 'en', 'Max workpiece size', '1000 × 800 × 600 mm'),
(1, 'pl', 'Maks. wymiar przedmiotu', '1000 × 800 × 600 mm'),
(2, 'ru', 'Мощность шпинделя', '15 кВт'),
(2, 'en', 'Spindle power', '15 kW'),
(2, 'pl', 'Moc wrzeciona', '15 kW'),
(3, 'ru', 'Точность позиционирования', '±0.005 мм'),
(3, 'en', 'Positioning accuracy', '±0.005 mm'),
(3, 'pl', 'Dokładność pozycjonowania', '±0.005 mm'),
(4, 'ru', 'Вес', '4500 кг'),
(4, 'en', 'Weight', '4500 kg'),
(4, 'pl', 'Waga', '4500 kg');

INSERT INTO news (image_url, published, created_at) VALUES
('/images/news/opening.jpg', true, NOW() - INTERVAL '5 days'),
('/images/news/exhibition.jpg', true, NOW() - INTERVAL '10 days'),
('/images/news/certificate.jpg', true, NOW() - INTERVAL '15 days');

INSERT INTO news_translations (news_id, language_code, title, content, excerpt) VALUES
(1, 'ru', 'Открытие нового цеха', 'Состоялось торжественное открытие нового производственного цеха площадью 5000 кв.м.', 'Новый цех увеличит производственные мощности'),
(1, 'en', 'New workshop opening', 'Grand opening of a new production workshop with an area of 5000 sq.m.', 'New workshop will increase production capacity'),
(1, 'pl', 'Otwarcie nowego warsztatu', 'Uroczyste otwarcie nowego warsztatu produkcyjnego o powierzchni 5000 mkw.', 'Nowy warsztat zwiększy moce produkcyjne'),
(2, 'ru', 'Участие в выставке', 'Принимаем участие в международной выставке промышленного оборудования.', 'Наш стенд на выставке IndustrialExpo 2024'),
(2, 'en', 'Exhibition participation', 'We are participating in the international industrial equipment exhibition.', 'Our booth at IndustrialExpo 2024'),
(2, 'pl', 'Udział w wystawie', 'Bierzemy udział w międzynarodowej wystawie sprzętu przemysłowego.', 'Nasze stoisko na IndustrialExpo 2024'),
(3, 'ru', 'Получен новый сертификат', 'Получили сертификат качества ISO 9001:2024.', 'Подтверждение качества нашей продукции'),
(3, 'en', 'New certificate received', 'Received ISO 9001:2024 quality certificate.', 'Confirmation of our product quality'),
(3, 'pl', 'Otrzymano nowy certyfikat', 'Otrzymano certyfikat jakości ISO 9001:2024.', 'Potwierdzenie jakości naszych produktów');

INSERT INTO documents (file_url, type, created_at) VALUES
('/files/gosts/gost-12345.pdf', 'gost', NOW()),
('/files/certificates/iso-9001.pdf', 'certificate', NOW()),
('/files/references/technical-guide.pdf', 'reference', NOW());

INSERT INTO document_translations (document_id, language_code, title, description) VALUES
(1, 'ru', 'ГОСТ 12345-2024', 'Стандарт на промышленное оборудование'),
(1, 'en', 'GOST 12345-2024', 'Standard for industrial equipment'),
(1, 'pl', 'GOST 12345-2024', 'Norma dotycząca sprzętu przemysłowego'),
(2, 'ru', 'Сертификат ISO 9001:2024', 'Система менеджмента качества'),
(2, 'en', 'ISO 9001:2024 Certificate', 'Quality management system'),
(2, 'pl', 'Certyfikat ISO 9001:2024', 'System zarządzania jakością'),
(3, 'ru', 'Технический справочник', 'Руководство по эксплуатации оборудования'),
(3, 'en', 'Technical Reference Guide', 'Equipment operation manual'),
(3, 'pl', 'Poradnik techniczny', 'Instrukcja obsługi sprzętu');
