    // Theme switching functionality
    const themeSwitcher = document.getElementById('themeSwitcher');
    const body = document.body;

    // Check saved theme in localStorage
    const currentTheme = localStorage.getItem('theme') || 'light';

    if (currentTheme === 'dark') {
    body.classList.add('dark-mode');
    themeSwitcher.textContent = 'Light Mode☀️';
}

    themeSwitcher.addEventListener('click', () => {
    body.classList.toggle('dark-mode');

    if (body.classList.contains('dark-mode')) {
    localStorage.setItem('theme', 'dark');
    themeSwitcher.textContent = 'Light Mode☀️';
} else {
    localStorage.setItem('theme', 'light');
    themeSwitcher.textContent = 'Dark Mode🌙';
}
});
    document.getElementById('loginForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    const errorElement = document.getElementById('errorMessage');
    errorElement.textContent = '';
    errorElement.classList.remove('show');

    try {
    const response = await fetch('/login', {
    method: 'POST',
    body: new FormData(this),
    headers: { 'Accept': 'application/json' },
    credentials: 'include' // Важно для работы с cookies
});

    // Проверяем Content-Type ответа
    const contentType = response.headers.get('content-type');
    let result;

    if (contentType && contentType.includes('application/json')) {
    result = await response.json();
} else {
    // Если ответ не JSON, обрабатываем как текст/HTML
    const text = await response.text();

    if (response.ok) {
    // Если статус 200, но не JSON - вероятно, это редирект через HTML
    window.location.href = '/protected/url';
    return;
}

    // Пытаемся найти ошибку в HTML (если это страница с ошибкой)
    const errorMatch = text.match(/<div class="error">(.*?)<\/div>/i);
    errorElement.textContent = errorMatch ? errorMatch[1] : 'Неизвестная ошибка сервера';
    errorElement.classList.add('show');
    return;
}

    if (!response.ok) {
    // Обработка JSON-ошибок от сервера
    errorElement.textContent = result.error || result.message || 'Ошибка авторизации';
    errorElement.classList.add('show');
    return;
}

    // Успешная авторизация
    window.location.href = result.redirect || '/protected/url"';

} catch (error) {
    console.error('Login error:', error);

    let errorMessage;
    if (error instanceof TypeError) {
    errorMessage = 'Нет соединения с сервером';
} else if (error instanceof SyntaxError) {
    errorMessage = 'Сервер вернул некорректные данные';
} else {
    errorMessage = 'Произошла неизвестная ошибка';
}

    errorElement.textContent = errorMessage;
    errorElement.classList.add('show');

    // Автоскрытие через 5 сек
    setTimeout(() => errorElement.classList.remove('show'), 5000);
}

});
