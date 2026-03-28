import en from './en.json';
import id from './id.json';

type TranslationKey = string;
type Translations = Record<string, any>;

const translations: Record<string, Translations> = {
    en,
    id
};

let currentLocale = 'en';

// Initialize locale from localStorage or browser preference
(function initLocale() {
    const savedLocale = localStorage.getItem('locale');
    if (savedLocale && translations[savedLocale]) {
        currentLocale = savedLocale;
    } else {
        const browserLang = navigator.language.split('-')[0];
        if (translations[browserLang]) {
            currentLocale = browserLang;
        }
    }
})();

/**
 * Get a nested value from an object using dot notation
 */
function getNestedValue(obj: Record<string, any>, path: string): any {
    return path.split('.').reduce((acc, part) => acc && acc[part], obj);
}

/**
 * Translate a key to the current locale
 * @param key - Translation key (e.g., 'auth.login')
 * @param params - Optional parameters to replace in the translation
 * @returns Translated string or the key if not found
 */
export function t(key: TranslationKey, params: Record<string, string | number> = {}): string {
    const translation = getNestedValue(translations[currentLocale], key);
    
    if (!translation) {
        // Fallback to English
        const fallback = getNestedValue(translations['en'], key);
        if (fallback) {
            return interpolate(fallback, params);
        }
        return key;
    }
    
    return interpolate(translation, params);
}

/**
 * Interpolate parameters into a translation string
 */
function interpolate(str: string, params: Record<string, string | number>): string {
    return str.replace(/\{(\w+)\}/g, (_, key) => {
        return params[key] !== undefined ? String(params[key]) : `{${key}}`;
    });
}

/**
 * Set the current locale
 */
export function setLocale(locale: string): void {
    if (translations[locale]) {
        currentLocale = locale;
        localStorage.setItem('locale', locale);
        document.documentElement.setAttribute('lang', locale);
    }
}

/**
 * Get the current locale
 */
export function getLocale(): string {
    return currentLocale;
}

/**
 * Get all available locales
 */
export function getAvailableLocales(): string[] {
    return Object.keys(translations);
}

// Export locale stores for Svelte reactivity
export const locales = Object.keys(translations);
