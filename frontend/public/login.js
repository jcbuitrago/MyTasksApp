// Espera a que el documento HTML esté completamente cargado
document.addEventListener('DOMContentLoaded', () => {

    const BASE_URL = 'http://localhost:8080';
    
    const loginForm = document.getElementById('loginForm');
    const usernameInput = document.getElementById('username');
    const passwordInput = document.getElementById('password');
    const errorMessageDiv = document.getElementById('errorMessage');
  
    loginForm.addEventListener('submit', async (event) => {
      event.preventDefault();
  
      errorMessageDiv.classList.add('d-none');
      
      const username = usernameInput.value.trim();
      const password = passwordInput.value.trim();
      const requestBody = {
        username: username,
        password: password
      };
      
      try {
        const response = await fetch(`${BASE_URL}/usuarios/iniciar-sesion`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(requestBody)
        });
  
        const data = await response.json();
  
        if (!response.ok) {
          throw new Error(data.message || 'Hubo un problema al iniciar sesión. Revisa tus credenciales.');
        }
  
        console.log('Login exitoso:', data);
        
        localStorage.setItem('authToken', data.token);
        localStorage.setItem('userName', data.user_name);
  
        window.location.href = 'index.html';
  
      } catch (error) {
        console.error('Error de login:', error);
        
        errorMessageDiv.textContent = error.message;
        errorMessageDiv.classList.remove('d-none');
      }
    });
  });