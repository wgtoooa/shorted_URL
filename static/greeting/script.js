
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
