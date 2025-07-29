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

// Password toggle visibility
const togglePassword = document.getElementById('togglePassword');
const passwordInput = document.getElementById('password');

togglePassword.addEventListener('click', () => {
    const type = passwordInput.getAttribute('type') === 'password' ? 'text' : 'password';
    passwordInput.setAttribute('type', type);
    togglePassword.textContent = type === 'password' ? '👁️' : '👁️‍🗨️';
});

// Form submission handling
document.getElementById('registerForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    const form = e.target;
    const errorEl = document.getElementById('errorMessage');
    errorEl.textContent = '';
    errorEl.classList.remove('show');

    try {
        const formData = new FormData(form);
        const response = await fetch('/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                login: formData.get('login'),
                password: formData.get('password')
            })
        });

        const result = await response.json();

        if (!response.ok) {
            throw new Error(result.error || 'Registration failed');
        }

        // Успешная регистрация
        window.location.href = result.redirect || '/login'; // Перенаправление клиентом

    } catch (error) {
        errorEl.textContent = error.message;
        errorEl.classList.add('show');

        // Автоскрытие ошибки через 5 сек

    }
});

// Input validation
document.getElementById('username').addEventListener('input', function() {
    if (this.value.length < 3 && this.value.length > 0) {
        this.setCustomValidity('Username must be at least 3 characters');
    } else {
        this.setCustomValidity('');
    }
});

document.getElementById('password').addEventListener('input', function() {
    if (this.value.length < 6 && this.value.length > 0) {
        this.setCustomValidity('Password must be at least 6 characters');
    } else {
        this.setCustomValidity('');
    }
});