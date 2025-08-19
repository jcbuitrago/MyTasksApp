document.addEventListener('DOMContentLoaded', () => {

    const BASE_URL = 'http://localhost:8080';
  
    const registerForm = document.getElementById('registerForm');
    const usernameInput = document.getElementById('username');
    const passwordInput = document.getElementById('password');
    const confirmPasswordInput = document.getElementById('confirmPassword');
    const errorMessageDiv = document.getElementById('errorMessage');
  
    registerForm.addEventListener('submit', async (event) => {
      event.preventDefault();
  
      errorMessageDiv.classList.add('d-none');
  
      const username = usernameInput.value.trim();
      const password = passwordInput.value.trim();
      const confirmPassword = confirmPasswordInput.value.trim();
  
      if (!username || !password || !confirmPassword) {
        showError('Todos los campos son obligatorios.');
        return;
      }
  
      if (password !== confirmPassword) {
        showError('Las contraseñas no coinciden.');
        return;
      }
  
      const requestBody = {
        username: username,
        password: password,
        picture_url: "https://png.pngtree.com/png-clipart/20210915/ourmid/pngtree-avatar-placeholder-abstract-white-blue-green-png-image_3918476.jpg"
      };
      
      try {
        const response = await fetch(`${BASE_URL}/usuarios`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(requestBody)
        });
  
        const data = await response.json();
  
        if (!response.ok) {
          throw new Error(data.message || 'No se pudo completar el registro.');
        }
        alert('¡Registro exitoso! Ahora serás redirigido para iniciar sesión.');
        window.location.href = 'login.html';
        localStorage.setItem('picture_url', data.picture_url);
  
      } catch (error) {
        showError(error.message);
      }
    });
    function showError(message) {
      errorMessageDiv.textContent = message;
      errorMessageDiv.classList.remove('d-none');
    }
  });