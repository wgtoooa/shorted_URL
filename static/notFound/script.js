    const themeSwitcher = document.getElementById('themeSwitcher');
    const body = document.body;

    // Проверяем сохранённую тему
    if (localStorage.getItem('darkMode') === 'enabled') {
    body.classList.add('dark-mode');
    themeSwitcher.textContent = 'Светлая тема';
}

    // Обработчик переключения темы
    themeSwitcher.addEventListener('click', () => {
    body.classList.toggle('dark-mode');

    if (body.classList.contains('dark-mode')) {
    localStorage.setItem('darkMode', 'enabled');
    themeSwitcher.textContent = 'Светлая тема';
} else {
    localStorage.setItem('darkMode', 'disabled');
    themeSwitcher.textContent = 'Тёмная тема';
}
});
