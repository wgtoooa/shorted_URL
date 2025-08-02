// Функция для экранирования HTML
function escapeHTML(str) {
    if (!str) return '';
    return str.toString()
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}

// Переключение темы
const themeSwitcher = document.getElementById('themeSwitcher');
const body = document.body;

// Проверяем сохранённую тему
const currentTheme = localStorage.getItem('theme') || 'light';

if (currentTheme === 'dark') {
    body.classList.add('dark-mode');
    themeSwitcher.textContent = 'Светлая тема☀️';
}

themeSwitcher.addEventListener('click', () => {
    body.classList.toggle('dark-mode');

    if (body.classList.contains('dark-mode')) {
        localStorage.setItem('theme', 'dark');
        themeSwitcher.textContent = 'Светлая тема☀️';
    } else {
        localStorage.setItem('theme', 'light');
        themeSwitcher.textContent = 'Темная тема🌙';
    }
});

// Модальное окно для удаления
const deleteModal = document.getElementById('deleteModal');
const modalTitle = document.getElementById('modalTitle');
const modalMessage = document.getElementById('modalMessage');
const cancelDeleteBtn = document.getElementById('cancelDelete');
const confirmDeleteBtn = document.getElementById('confirmDelete');
const deleteMessage = document.getElementById('deleteMessage');

let currentShortUrlToDelete = null;

function showModal(title, message) {
    modalTitle.textContent = title;
    modalMessage.textContent = message;
    deleteMessage.textContent = '';
    deleteMessage.className = 'message';
    deleteModal.style.display = 'block';
}

function hideModal() {
    deleteModal.style.display = 'none';
    currentShortUrlToDelete = null;
}

// Добавляем обработчик для кнопки подтверждения удаления
confirmDeleteBtn.addEventListener('click', () => {
    if (currentShortUrlToDelete) {
        deleteLink(currentShortUrlToDelete);
    }
});

cancelDeleteBtn.addEventListener('click', hideModal);

window.addEventListener('click', (event) => {
    if (event.target === deleteModal) {
        hideModal();
    }
});

async function deleteLink(shortUrl) {
    try {
        confirmDeleteBtn.disabled = true;
        confirmDeleteBtn.textContent = 'Удаление...';

        const res = await fetch(`/protected/url/${encodeURIComponent(shortUrl)}`, {
            method: 'DELETE',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (!res.ok) {
            const errorData = await res.json().catch(() => ({}));
            throw new Error(errorData.message || 'Ошибка при удалении ссылки');
        }

        deleteMessage.textContent = 'Ссылка успешно удалена';
        deleteMessage.className = 'message success';

        setTimeout(() => {
            loadLinks();
            hideModal();
        }, 1000);
    } catch (err) {
        console.error('Delete error:', err);
        deleteMessage.textContent = 'Ошибка: ' + (err.message || 'Неизвестная ошибка');
        deleteMessage.className = 'message error';
    } finally {
        confirmDeleteBtn.disabled = false;
        confirmDeleteBtn.textContent = 'Удалить';
    }
}

function showDeleteModal(shortUrl) {
    currentShortUrlToDelete = shortUrl;
    showModal("Подтвердите удаление", `Вы уверены, что хотите удалить ссылку "${shortUrl}"?`);
}

// Загрузка и отображение списка ссылок
async function loadLinks() {
    const container = document.getElementById('urlList');
    if (!container) return;

    container.innerHTML = '<div class="loading">Загрузка ссылок...</div>';

    try {
        const res = await fetch('/protected/url/data', {
            method: 'GET',
            credentials: 'include'
        });

        if (res.status === 401) {
            window.location.href = '/login';
            return;
        }

        if (!res.ok) {
            throw new Error(`Ошибка HTTP: ${res.status}`);
        }

        const data = await res.json();
        console.log('Data received:', data);

        container.innerHTML = '';

        if (!data || !Array.isArray(data) || data.length === 0) {
            container.innerHTML = '<div class="empty">У вас ещё нет ссылок.</div>';
            return;
        }

        data.forEach(link => {
            if (!link.short_url || !link.full_url) {
                console.warn('Invalid link format:', link);
                return;
            }

            const div = document.createElement('div');
            div.className = 'url-item';
            div.innerHTML = `
                <a class="short-url" href="l/${escapeHTML(link.short_url)}" target="_blank">
                    ${escapeHTML(link.short_url)}
                </a>
                <div class="original-url">${escapeHTML(link.full_url)}</div>
                <div class="dropdown">
                    <button class="dropdown-toggle">⋯</button>
                    <div class="dropdown-menu">
                        <a href="#" class="dropdown-item edit-item" data-short-url="${escapeHTML(link.short_url)}">
                            Изменить короткую ссылку
                        </a>
                        <a href="#" class="dropdown-item copy-item" data-full-url="${escapeHTML(link.full_url)}">
                            Копировать
                        </a>
                        <div class="dropdown-divider"></div>
                        <a href="#" class="dropdown-item delete-item" data-short-url="${escapeHTML(link.short_url)}">
                            Удалить
                        </a>
                    </div>
                </div>
            `;
            container.appendChild(div);
        });

        // Инициализация событий
        initDropdowns();
        initCopyButtons();
        initDeleteButtons();
        initEditButtons();
        initShortUrlClicks();

    } catch (err) {
        console.error('Error loading links:', err);
        container.innerHTML = `
            <div class="error">
                Ошибка загрузки: ${err.message}
                <button onclick="loadLinks()">Повторить</button>
            </div>
        `;
    }
}

// Инициализация кликов по коротким ссылкам
function initShortUrlClicks() {
    document.querySelectorAll('.short-url').forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            const url = this.getAttribute('href');
            window.location.href = url;
        });
    });
}

// Инициализация dropdown-меню
function initDropdowns() {
    document.querySelectorAll('.dropdown-toggle').forEach(toggle => {
        toggle.addEventListener('click', function(e) {
            e.preventDefault();
            const menu = this.nextElementSibling;
            document.querySelectorAll('.dropdown-menu').forEach(m => {
                if (m !== menu) m.classList.remove('show');
            });
            menu.classList.toggle('show');
        });
    });

    document.addEventListener('click', function(e) {
        if (!e.target.matches('.dropdown-toggle') && !e.target.closest('.dropdown-menu')) {
            document.querySelectorAll('.dropdown-menu').forEach(menu => {
                menu.classList.remove('show');
            });
        }
    });
}

// Инициализация кнопок копирования
function initCopyButtons() {
    document.querySelectorAll('.copy-item').forEach(item => {
        item.addEventListener('click', async function(e) {
            e.preventDefault();
            const fullUrl = this.dataset.fullUrl;
            try {
                await navigator.clipboard.writeText(fullUrl);
                showToast('Ссылка скопирована!');
            } catch (err) {
                console.error('Ошибка копирования:', err);
                showToast('Не удалось скопировать ссылку', 'error');
            }
        });
    });
}

// Инициализация кнопок удаления
function initDeleteButtons() {
    document.querySelectorAll('.delete-item').forEach(item => {
        item.addEventListener('click', function(e) {
            e.preventDefault();
            const shortUrl = this.dataset.shortUrl;
            showDeleteModal(shortUrl);
        });
    });
}
// Основная функция инициализации
function initEditButtons() {
    document.querySelectorAll('.edit-item').forEach(button => {
        button.addEventListener('click', handleEditClick);
    });
}

// Обработчик клика
function handleEditClick(e) {
    e.preventDefault();
    const { shortUrl } = this.dataset;

    if (!shortUrl) {
        showToast('Недостаточно данных для редактирования', 'error');
        return;
    }

    const modal = createEditModal(shortUrl);
    document.body.appendChild(modal);
    setupModalEvents(modal, shortUrl);
}

// Создание модального окна
function createEditModal(shortUrl) {
    const modal = document.createElement('div');
    modal.className = 'modal active';
    modal.innerHTML = `
        <div class="modal-content">
            <span class="close">&times;</span>
            <h3>Редактирование ссылки</h3>
            <div class="form-group">
                <label>Короткий URL:</label>
                <input type="text" id="short-url" value="${escapeHTML(shortUrl)}">
                <div id="error-message" class="error-message"></div>
            </div>
            <div class="modal-buttons">
                <button class="modal-button cancel">Отмена</button>
                <button id="save-btn" class="modal-button confirm">Сохранить</button>
            </div>
        </div>
    `;
    return modal;
}

// Настройка обработчиков событий
function setupModalEvents(modal, originalUrl) {
    const closeModal = () => modal.remove();

    // Обработчики закрытия
    modal.addEventListener('click', e => {
        if (e.target === modal || e.target.classList.contains('close') || e.target.classList.contains('cancel')) {
            closeModal();
        }
    });

    // Обработчик сохранения
    modal.querySelector('#save-btn').addEventListener('click', async () => {
        await handleSaveAction(modal, originalUrl);
    });
}

// Обработчик сохранения
async function handleSaveAction(modal, originalUrl) {
    const newUrl = modal.querySelector('#short-url').value.trim();
    const errorElement = modal.querySelector('#error-message');
    const saveBtn = modal.querySelector('#save-btn');

    // Сброс предыдущих ошибок
    errorElement.textContent = '';

    // Валидация
    if (!newUrl) {
        errorElement.textContent = 'Поле не может быть пустым';
        return;
    }

    if (newUrl === originalUrl) {
        showToast('Изменений не обнаружено', 'info');
        modal.remove();
        return;
    }

    try {
        // Блокировка кнопки на время запроса
        saveBtn.disabled = true;
        saveBtn.textContent = 'Сохранение...';

        const response = await fetch('/protected/url', {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: JSON.stringify({
                old_short_url: originalUrl,
                new_short_url: newUrl
            }),
            credentials: 'include'
        });

        await processResponse(response, modal, errorElement);
    } catch (error) {
        handleNetworkError(error, errorElement);
    } finally {
        saveBtn.disabled = false;
        saveBtn.textContent = 'Сохранить';
    }
}

// Обработка ответа сервера
async function processResponse(response, modal, errorElement) {
    if (response.ok) {
        const data = await parseJSON(response);
        handleSuccess(modal, data?.message);
    } else {
        await handleServerError(response, errorElement);
    }
}

// Успешное выполнение
function handleSuccess(modal, message) {
    modal.remove();
    showToast(message || 'Ссылка успешно обновлена', 'success');
    setTimeout(() => location.reload(), 1000);
}

// Ошибки сервера
async function handleServerError(response, errorElement) {
    try {
        const errorData = await parseJSON(response);
        errorElement.textContent = errorData?.message ||
            `Ошибка ${response.status}: ${response.statusText}`;
    } catch {
        errorElement.textContent = `Ошибка ${response.status}: ${response.statusText}`;
    }
}

// Сетевые ошибки
function handleNetworkError(error, errorElement) {
    console.error('Ошибка сети:', error);
    errorElement.textContent = 'Ошибка соединения с сервером';
    showToast('Не удалось подключиться к серверу', 'error');
}

// Парсинг JSON с защитой
async function parseJSON(response) {
    const text = await response.text();
    try {
        return text ? JSON.parse(text) : {};
    } catch {
        return { message: text };
    }
}

// Ваша функция экранирования
function escapeHTML(str) {
    if (!str) return '';
    return str.toString()
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}

// Инициализация после загрузки DOM
document.addEventListener('DOMContentLoaded', initEditButtons);
// Функция для показа уведомлений
function showToast(message, type = 'success') {
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;
    document.body.appendChild(toast);

    setTimeout(() => {
        toast.classList.add('show');
    }, 10);

    setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => {
            toast.remove();
        }, 300);
    }, 3000);
}

// Обработка формы создания ссылки
document.getElementById('createForm')?.addEventListener('submit', async function(e) {
    e.preventDefault();
    const form = e.target;
    const full_url = form.elements.full_url.value.trim();
    const short_url = form.elements.short_url?.value.trim();
    const messageEl = document.getElementById('createMessage');

    messageEl.textContent = '';
    messageEl.classList.remove('error', 'success');

    // Валидация URL
    if (!/^https?:\/\//i.test(full_url)) {
        messageEl.textContent = 'URL должен начинаться с http:// или https://';
        messageEl.classList.add('error');
        return;
    }

    // Валидация короткой ссылки (если указана)
    // if (short_url && !/^[a-zA-Z0-9_-]+$/.test(short_url)) {
    //     messageEl.textContent = 'Короткая ссылка может содержать только буквы, цифры, "_" и "-"';
    //     messageEl.classList.add('error');
    //     return;
    // }

    const data = { full_url };
    if (short_url) data.short_url = short_url;

    try {
        const submitBtn = form.querySelector('button[type="submit"]');
        submitBtn.disabled = true;

        const res = await fetch('/protected/url', {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        const result = await res.json();

        if (!res.ok) {
            throw new Error(result.error || 'Ошибка при создании ссылки');
        }

        messageEl.textContent = 'Ссылка успешно создана!';
        messageEl.classList.add('success');
        form.reset();
        loadLinks();
    } catch (err) {
        messageEl.textContent = 'Ошибка: ' + err.message;
        messageEl.classList.add('error');
    } finally {
        const submitBtn = form.querySelector('button[type="submit"]');
        if (submitBtn) submitBtn.disabled = false;
    }
});

// Загрузка при старте
window.addEventListener('DOMContentLoaded', () => {
    loadLinks();

    // Инициализация всех элементов
    if (themeSwitcher) {
        const currentTheme = localStorage.getItem('theme') || 'light';
        if (currentTheme === 'dark') {
            body.classList.add('dark-mode');
        }
    }
});