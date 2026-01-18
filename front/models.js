class Page {
    constructor(data) {
        this.id = data.page_id || data.id;
        this.slug = data.slug;
        this.template = data.template;
        this.translations = data.translations || [];
        this.createdAt = data.created_at;
        this.updatedAt = data.updated_at;
    }
    
    getTranslation(locale) {
        const translation = this.translations.find(t => t.language_code === locale);
        if (translation) return translation;
        
        const ruTranslation = this.translations.find(t => t.language_code === 'ru');
        if (ruTranslation) return ruTranslation;
        
        return {};
    }
    
    getTranslationsMap() {
        const map = {};
        this.translations.forEach(translation => {
            map[translation.language_code] = translation;
        });
        return map;
    }
}

class ProductCategory {
    constructor(data) {
        this.id = data.id;
        this.parentId = data.parent_id;
        this.sortOrder = data.sort_order;
        this.translationsArray = data.translations || [];
        this.translations = this.createTranslationsMap();
    }
    
    createTranslationsMap() {
        const map = {};
        this.translationsArray.forEach(translation => {
            map[translation.language_code] = translation;
        });
        return map;
    }
    
    getName(locale) {
        const translation = this.translations[locale] || this.translations['ru'];
        return translation ? translation.name : '';
    }
}

class Product {
    constructor(data) {
        this.id = data.id;
        this.categoryId = data.category_id;
        this.sku = data.sku;
        this.imageUrl = data.image_url;
        this.fileUrl = data.file_url;
        this.sortOrder = data.sort_order;
        this.translationsArray = data.translations || [];
        this.translations = this.createTranslationsMap();
        this.specs = data.specs || [];
        this.category = data.category;
    }
    
    createTranslationsMap() {
        const map = {};
        this.translationsArray.forEach(translation => {
            map[translation.language_code] = translation;
        });
        return map;
    }
    
    getName(locale) {
        const translation = this.translations[locale] || this.translations['ru'];
        return translation ? translation.name : '';
    }
    
    getDescription(locale) {
        const translation = this.translations[locale] || this.translations['ru'];
        return translation ? translation.description : '';
    }
}

class News {
    constructor(data) {
        this.id = data.id;
        this.imageUrl = data.image_url;
        this.published = data.published;
        this.createdAt = data.created_at;
        this.updatedAt = data.updated_at;
        this.translations = this.createTranslationsMap;
    }

     createTranslationsMap() {
        const map = {};
        this.translationsArray.forEach(translation => {
            map[translation.language_code] = translation;
        });
        return map;
    }
    
    getTitle(locale) {
        const translation = this.translations[locale] || this.translations['ru'];
        return translation ? translation.title : '';
    }
    
    getContent(locale) {
        const translation = this.translations[locale] || this.translations['ru'];
        return translation ? translation.content : '';
    }
}

class Document {
    constructor(data) {
        this.id = data.id;
        this.fileUrl = data.file_url;
        this.type = data.type;
        this.createdAt = data.created_at;
        this.translations = this.createTranslationsMap;
    }
    
    createTranslationsMap() {
        const map = {};
        this.translationsArray.forEach(translation => {
            map[translation.language_code] = translation;
        });
        return map;
    }
    
    getTitle(locale) {
        const translation = this.translations[locale] || this.translations['ru'];
        return translation ? translation.title : '';
    }
}

class Contact {
    constructor(data) {
        this.id = data.id;
        this.type = data.type;
        this.value = data.value;
        this.sortOrder = data.sort_order;
        this.translations = this.createTranslationsMap;
    }
    
    createTranslationsMap() {
        const map = {};
        this.translationsArray.forEach(translation => {
            map[translation.language_code] = translation;
        });
        return map;
    }
    
    getLabel(locale) {
        const translation = this.translations[locale] || this.translations['ru'];
        return translation ? translation.label : '';
    }
}

class Feedback {
    constructor(data) {
        this.name = data.name;
        this.email = data.email;
        this.phone = data.phone;
        this.company = data.company;
        this.message = data.message;
    }
}

class Store {
    constructor() {
        this.state = {
            locale: 'ru',
            currentPage: null,
            categories: [],
            products: [],
            news: [],
            documents: [],
            contacts: [],
            searchQuery: '',
            isLoading: false
        };
        
        this.observers = [];
    }
    
    setState(newState) {
        this.state = { ...this.state, ...newState };
        this.notifyObservers();
    }
    
    subscribe(observer) {
        this.observers.push(observer);
    }
    
    notifyObservers() {
        this.observers.forEach(observer => observer(this.state));
    }
    
    changeLocale(locale) {
        this.setState({ locale });
        localStorage.setItem('locale', locale);
        document.documentElement.lang = locale;
    }

    findNews(newsID){

    }

    findProduct(productID){
        
    }
}

window.Models = {
    Page,
    ProductCategory,
    Product,
    News,
    Document,
    Contact,
    Feedback,
    Store
};