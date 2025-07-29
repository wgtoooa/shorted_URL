const themeSwitcher = document.getElementById('themeSwitcher');
const body = document.body;

// Проверяем сохранённую тему в localStorage
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

// Функция для показа сообщения о копировании
function showCopyMessage() {
    const message = document.getElementById('copyMessage');
    message.style.display = 'block';
    setTimeout(() => {
        message.style.display = 'none';
    }, 2000);
}

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

cancelDeleteBtn.addEventListener('click', hideModal);

// Закрытие модального окна при клике вне его
window.addEventListener('click', (event) => {
    if (event.target === deleteModal) {
        hideModal();
    }
});

async function deleteLink(shortUrl) {
    try {
        const res = await fetch(`/url/${encodeURIComponent(shortUrl)}`, {
            method: 'DELETE',
            credentials: 'include'
        });

        if (!res.ok) {
            throw new Error('Ошибка при удалении ссылки');
        }

        deleteMessage.textContent = 'Ссылка успешно удалена';
        deleteMessage.className = 'message';

        // Обновляем список через 1 секунду, чтобы пользователь увидел сообщение
        setTimeout(() => {
            loadLinks();
            hideModal();
        }, 1000);
    } catch (err) {
        deleteMessage.textContent = 'Ошибка: ' + err.message;
        deleteMessage.className = 'message error';
    }
}

confirmDeleteBtn.addEventListener('click', () => {
    if (currentShortUrlToDelete) {
        deleteLink(currentShortUrlToDelete);
    }
});


// ... (всё как раньше, до loadLinks)

function escapeHTML(str) {
    return str.replace(/[&<>"']/g, function (m) {
        return {
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            '"': '&quot;',
            "'": '&#39;'
        }[m];
    });
}

function showDeleteModal(shortUrl) {
    currentShortUrlToDelete = shortUrl;
    showModal("Подтвердите удаление", "Вы уверены, что хотите удалить эту ссылку?");
}

async function loadLinks() {
    try {
        const res = await fetch('/url/data', {
            method: 'GET',
            credentials: 'include'
        });

        const container = document.getElementById('urlList');

        if (!res.ok) {
            container.textContent = 'Ошибка при загрузке ссылок';
            return;
        }

        const urls = await res.json();
        container.innerHTML = '';

        if (!urls || urls.length === 0) {
            container.innerHTML = 'У вас ещё нет ссылок.';
            return;
        }

        urls.forEach(link => {
            const div = document.createElement('div');
            div.className = 'url-item';
            div.innerHTML = `
                    <a class="short-url" href="${escapeHTML(link.short_url)}" target="_blank">${escapeHTML(link.short_url)}</a>
                    <div class="original-url">${escapeHTML(link.full_url)}</div>
                    <div class="dropdown">
                        <button class="dropdown-toggle">⋯</button>
                        <div class="dropdown-menu">
                            <a href="#" class="dropdown-item edit-item" data-short-url="${escapeHTML(link.short_url)}">Изменить</a>
                            <a href="#" class="dropdown-item copy-item" data-full-url="${escapeHTML(link.full_url)}">Копировать полную ссылку</a>
                            <div class="dropdown-divider"></div>
                            <a href="#" class="dropdown-item delete-item" data-short-url="${escapeHTML(link.short_url)}">Удалить</a>
                        </div>
                    </div>
                `;
            container.appendChild(div);
        });

        if (navigator.clipboard && navigator.clipboard.writeText) {
            navigator.clipboard.writeText(fullUrl)
                .then(() => {
                    showCopyMessage();
                })
                .catch(err => {
                    console.error('Ошибка при копировании: ', err.message);
                    alert('Не удалось скопировать ссылку: ' + err.message);
                });
        } else {
            alert("Копирование не поддерживается в этом браузере.");
        }


        document.querySelectorAll('.dropdown-toggle').forEach(toggle => {
            toggle.addEventListener('click', function (e) {
                e.preventDefault();
                const menu = this.nextElementSibling;
                document.querySelectorAll('.dropdown-menu').forEach(m => {
                    if (m !== menu) m.classList.remove('show');
                });
                menu.classList.toggle('show');
            });
        });

        document.addEventListener('click', function (e) {
            if (!e.target.closest('.dropdown')) {
                document.querySelectorAll('.dropdown-menu').forEach(menu => {
                    menu.classList.remove('show');
                });
            }
        });

        document.querySelectorAll('.delete-item').forEach(item => {
            item.addEventListener('click', function (e) {
                e.preventDefault();
                const shortUrl = this.dataset.shortUrl;
                showDeleteModal(shortUrl);
            });
        });



        document.querySelectorAll('.edit-item').forEach(item => {
            item.addEventListener('click', function (e) {
                e.preventDefault();
                const shortUrl = this.dataset.shortUrl;
                alert(`Функция изменения ссылки "${shortUrl}" будет реализована позже.`);
            });
        });

    } catch (err) {
        document.getElementById('urlList').textContent = 'Ошибка: ' + err.message;
    }
}

document.getElementById('createForm').addEventListener('submit', async function (e) {
    e.preventDefault();
    const full_url = document.getElementById('full_url').value.trim();
    const short_url = document.getElementById('short_url').value.trim();
    const messageEl = document.getElementById('createMessage');

    messageEl.textContent = '';
    messageEl.classList.remove('error');

    if (!/^https?:\/\//.test(full_url)) {
        messageEl.textContent = 'Ссылка должна начинаться с http:// или https://';
        messageEl.classList.add('error');
        return;
    }

    const data = { full_url };
    if (short_url) data.short_url = short_url;

    try {
        const res = await fetch('/url', {
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
        document.getElementById('createForm').reset();
        loadLinks();
    } catch (err) {
        messageEl.textContent = 'Ошибка: ' + err.message;
        messageEl.classList.add('error');
    }
});

window.addEventListener('DOMContentLoaded', loadLinks);