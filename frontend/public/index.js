// --- GUARDIA DE SESIÓN ---
const authToken = localStorage.getItem('authToken');
if (!authToken) {
  window.location.href = 'login.html';
}

// --- LÓGICA DEL DASHBOARD ---
document.addEventListener('DOMContentLoaded', () => {

  const BASE_URL = 'http://localhost:8080';

  // --- Selección de Elementos del DOM (una sola vez) ---
  const taskListContainer = document.getElementById('taskListContainer');
  const logoutButton = document.getElementById('logoutButton');

  // Elementos del modal de AÑADIR tarea
  const addModal = new bootstrap.Modal(document.getElementById('addTaskModal'));
  const addTaskForm = document.getElementById('addTaskForm');
  const addTaskErrorDiv = document.getElementById('addTaskError');

  // Elementos del modal de EDITAR tarea
  const editModal = new bootstrap.Modal(document.getElementById('editTaskModal'));
  const editTaskForm = document.getElementById('editTaskForm');
  const editTaskErrorDiv = document.getElementById('editTaskError');
  const editTaskDeleteButton = document.getElementById('editTaskDelete'); // Asumiendo que tienes un botón de eliminar

  // --- Estado de la Aplicación (fuente única de verdad) ---
  let allTasks = [];
  let allCategories = [];
  let currentFilter = 'Sin Empezar';

  /**
   * 1. OBTIENE DATOS DE LA API
   * Funciones para cargar tareas y categorías UNA SOLA VEZ.
   */
  const fetchAllCategories = async () => {
    try {
      const response = await fetch(`${BASE_URL}/categorias`, {
        headers: { 'Authorization': `Bearer ${authToken}` }
      });
      if (!response.ok) throw new Error('No se pudieron cargar las categorías.');
      allCategories = await response.json(); // Guardamos las categorías en la variable global
    } catch (error) {
      console.error(error.message);
      // Si falla algo crítico, cerramos sesión
      localStorage.removeItem('authToken');
      window.location.href = 'login.html';
    }
  };

  const fetchAllTasks = async () => {
    try {
      const response = await fetch(`${BASE_URL}/tareas/usuario`, {
        headers: { 'Authorization': `Bearer ${authToken}` }
      });
      if (!response.ok) throw new Error('Error al obtener las tareas.');
      allTasks = await response.json();
      renderTasks(); // Una vez que tenemos las tareas, las mostramos
    } catch (error) {
      console.error(error.message);
      localStorage.removeItem('authToken');
      window.location.href = 'login.html';
    }
  };

  /**
   * 2. ACTUALIZA LA INTERFAZ (UI)
   * Funciones que dibujan los datos en la pantalla.
   */
  const populateCategoryDropdowns = () => {
    const selects = [
      document.getElementById('taskCategoryId'),
      document.getElementById('editTaskCategoryId')
    ];
    selects.forEach(selectElement => {
      if (selectElement) {
        selectElement.innerHTML = '<option value="">Selecciona una categoría</option>';
        allCategories.forEach(category => {
          const option = document.createElement('option');
          option.value = category.id;
          option.textContent = category.name;
          selectElement.appendChild(option);
        });
      }
    });
  };

  const renderTasks = () => {
    taskListContainer.innerHTML = '';
    let filteredTasks = [];

    if (currentFilter === 'Today') {
      const today = new Date().toISOString().split('T')[0];
      filteredTasks = allTasks.filter(task => task.tentative_due_date === today);
    } else {
      filteredTasks = allTasks.filter(task => task.status === currentFilter);
    }

    if (filteredTasks.length === 0) {
      taskListContainer.innerHTML = `<p class="text-center text-custom-blue p-4">No hay tareas en la sección "${currentFilter}".</p>`;
      return;
    }

    filteredTasks.forEach(task => {
      // **AJUSTE PRINCIPAL**: Buscamos el nombre de la categoría en el array que ya cargamos.
      const category = allCategories.find(c => c.id === task.category_id);
      const categoryName = category ? category.name : 'Sin categoría'; // Si no la encuentra, pone un texto por defecto

      const taskElement = document.createElement('div');
      taskElement.className = 'list-group-item d-flex align-items-center gap-4 bg-custom-dark-blue px-4 py-3 border-bottom border-custom-dark';
      taskElement.setAttribute('data-task-id', task.id);
      taskElement.setAttribute('data-bs-toggle', 'modal');
      taskElement.setAttribute('data-bs-target', '#editTaskModal');
      taskElement.style.cursor = 'pointer';
      
      taskElement.innerHTML = `
        <div class="d-flex flex-column justify-content-center">
            <p class="text-white fw-medium m-0">${task.description}</p>
            <p class="text-custom-blue small m-0">Vence: ${task.tentative_due_date}</p>
            <p class="text-info small m-0 fst-italic">Categoría: ${categoryName}</p>
        </div>
      `;
      taskListContainer.appendChild(taskElement);
    });
  };

  /**
   * 3. MANEJO DE EVENTOS DE USUARIO
   * Lógica para los clics en botones y formularios.
   */

  // Abrir y poblar el modal de edición
  taskListContainer.addEventListener('click', (event) => {
    const taskElement = event.target.closest('.list-group-item');
    if (taskElement) {
      const taskId = taskElement.dataset.taskId;
      const task = allTasks.find(t => t.id == taskId);
      if (task) {
        document.getElementById('editTaskId').value = task.id;
        document.getElementById('editTaskDescription').value = task.description;
        document.getElementById('editTaskDueDate').value = task.tentative_due_date;
        document.getElementById('editTaskStatus').value = task.status;
        document.getElementById('editTaskCategoryId').value = task.category_id;
      }
    }
  });

  // Enviar formulario para editar una tarea
  editTaskForm.addEventListener('submit', async (event) => {
    event.preventDefault();
    const id = document.getElementById('editTaskId').value;
    const requestBody = {
      description: document.getElementById('editTaskDescription').value.trim(),
      tentative_due_date: document.getElementById('editTaskDueDate').value,
      status: document.getElementById('editTaskStatus').value,
      category_id: parseInt(document.getElementById('editTaskCategoryId').value)
    };
    try {
      const response = await fetch(`${BASE_URL}/tareas/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${authToken}` },
        body: JSON.stringify(requestBody)
      });
      if (!response.ok) throw new Error('No se pudo actualizar la tarea.');
      editModal.hide();
      await fetchAllTasks(); // Recarga las tareas para reflejar el cambio
    } catch (error) {
      editTaskErrorDiv.textContent = error.message;
      editTaskErrorDiv.classList.remove('d-none');
    }
  });
  
  // Enviar formulario para añadir una tarea
  addTaskForm.addEventListener('submit', async (event) => {
    event.preventDefault();
    let requestBody;
    if (document.getElementById('taskCategoryId').value === '') {
      requestBody = {
        description: document.getElementById('taskDescription').value.trim(),
        tentative_due_date: document.getElementById('taskDueDate').value,
        status: document.getElementById('taskStatus').value,
      };
    } else {
      requestBody = {
        description: document.getElementById('taskDescription').value.trim(),
        tentative_due_date: document.getElementById('taskDueDate').value,
        status: document.getElementById('taskStatus').value,
        category_id: parseInt(document.getElementById('taskCategoryId').value)
      };
    }
    if (!requestBody.description || !requestBody.tentative_due_date) {
      addTaskErrorDiv.textContent = 'Por favor, completa todos los campos.';
      addTaskErrorDiv.classList.remove('d-none');
      return;
    }
    try {
      const response = await fetch(`${BASE_URL}/tareas`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${authToken}` },
        body: JSON.stringify(requestBody)
      });
      if (!response.ok) throw new Error('No se pudo crear la tarea.');
      addModal.hide();
      addTaskForm.reset();
      await fetchAllTasks(); // Recarga las tareas para mostrar la nueva
    } catch (error) {
      addTaskErrorDiv.textContent = error.message;
      addTaskErrorDiv.classList.remove('d-none');
    }
  });

  // Botón de Cerrar Sesión
  if (logoutButton) {
    logoutButton.addEventListener('click', () => {
      localStorage.removeItem('authToken');
      window.location.href = 'login.html';
    });
  }
  
  // Botones de Filtro
  document.querySelectorAll('[data-filter]').forEach(button => {
    button.addEventListener('click', (event) => {
      currentFilter = event.currentTarget.dataset.filter;
      renderTasks(); // Solo re-renderiza con el nuevo filtro, sin llamar a la API
    });
  });

  /**
   * 4. INICIALIZACIÓN
   * Orquesta la carga inicial de la aplicación.
   */
  const initializeApp = async () => {
    await fetchAllCategories();    // 1. Carga las categorías
    populateCategoryDropdowns(); // 2. Rellena los menús desplegables
    await fetchAllTasks();         // 3. Carga las tareas
  };

  editTaskDeleteButton.addEventListener('click', async () => {
    const id = document.getElementById('editTaskId').value;
    if (confirm('¿Estás seguro de que quieres eliminar esta tarea?')) {
      try {
        const response = await fetch(`${BASE_URL}/tareas/${id}`, {
          method: 'DELETE',
          headers: { 'Authorization': `Bearer ${authToken}` }
        });
        if (!response.ok) throw new Error('No se pudo eliminar la tarea.');
        editModal.hide();
        await fetchAllTasks();
      } catch (error) {
        const errorDiv = document.getElementById('editTaskError');
        errorDiv.textContent = error.message;
        errorDiv.classList.remove('d-none');
      }
    }
  });

  initializeApp();
});