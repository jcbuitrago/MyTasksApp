const authToken = localStorage.getItem('authToken');
if (!authToken) {
  window.location.href = 'login.html';
}

document.addEventListener('DOMContentLoaded', () => {

  const BASE_URL = 'http://localhost:8080';

  const taskListContainer = document.getElementById('taskListContainer');
  const logoutButton = document.getElementById('logoutButton');

  const addModal = new bootstrap.Modal(document.getElementById('addTaskModal'));
  const addTaskForm = document.getElementById('addTaskForm');
  const addTaskErrorDiv = document.getElementById('addTaskError');

  const editModal = new bootstrap.Modal(document.getElementById('editTaskModal'));
  const editTaskForm = document.getElementById('editTaskForm');
  const editTaskErrorDiv = document.getElementById('editTaskError');
  const editTaskDeleteButton = document.getElementById('editTaskDelete');

  const categoryFilterButton = document.getElementById('categoryFilterButton');
  const categoryFilterDropdown = document.getElementById('categoryFilterDropdown');

  let allTasks = [];
  let allCategories = [];
  let currentFilter = 'Sin Empezar';
  let currentCategoryFilter = 'all';

  const fetchAllCategories = async () => {
    try {
      const response = await fetch(`${BASE_URL}/categorias`, {
        headers: { 'Authorization': `Bearer ${authToken}` }
      });
      if (!response.ok) throw new Error('No se pudieron cargar las categorías.');
      allCategories = await response.json();
    } catch (error) {
      console.error(error.message);
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
      renderTasks();
    } catch (error) {
      console.error(error.message);
      localStorage.removeItem('authToken');
      window.location.href = 'login.html';
    }
  };

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

  const populateCategoryFilter = () => {
    if (categoryFilterDropdown) {
        categoryFilterDropdown.innerHTML = '';
        const allOption = document.createElement('li');
        allOption.innerHTML = `<a class="dropdown-item" href="#" data-category-id="all">Todas las Categorías</a>`;
        categoryFilterDropdown.appendChild(allOption);
        allCategories.forEach(category => {
            const option = document.createElement('li');
            option.innerHTML = `<a class="dropdown-item" href="#" data-category-id="${category.id}">${category.name}</a>`;
            categoryFilterDropdown.appendChild(option);
        });
    }
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

    if (currentCategoryFilter !== 'all') {
      filteredTasks = filteredTasks.filter(task => task.category_id == currentCategoryFilter);
    }

    if (filteredTasks.length === 0) {
      taskListContainer.innerHTML = `<p class="text-center text-custom-blue p-4">No hay tareas que coincidan con ${currentFilter}.</p>`;
      return;
    }

    filteredTasks.forEach(task => {
      const category = allCategories.find(c => c.id === task.category_id);
      const categoryName = category ? category.name : 'Sin categoría';
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
      await fetchAllTasks();
    } catch (error) {
      editTaskErrorDiv.textContent = error.message;
      editTaskErrorDiv.classList.remove('d-none');
    }
  });
  
  addTaskForm.addEventListener('submit', async (event) => {
    event.preventDefault();
    addTaskErrorDiv.classList.add('d-none');
    
    let requestBody;
    const description = document.getElementById('taskDescription').value.trim();
    const tentative_due_date = document.getElementById('taskDueDate').value;
    const status = document.getElementById('taskStatus').value;
    const categoryIdValue = document.getElementById('taskCategoryId').value;

    if (!description || !tentative_due_date) {
      addTaskErrorDiv.textContent = 'La descripción y la fecha son obligatorias.';
      addTaskErrorDiv.classList.remove('d-none');
      return;
    }

    if (categoryIdValue) {
        requestBody = { description, tentative_due_date, status, category_id: parseInt(categoryIdValue) };
    } else {
        requestBody = { description, tentative_due_date, status };
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
      await fetchAllTasks();
    } catch (error) {
      addTaskErrorDiv.textContent = error.message;
      addTaskErrorDiv.classList.remove('d-none');
    }
  });

  if (logoutButton) {
    logoutButton.addEventListener('click', () => {
      localStorage.removeItem('authToken');
      window.location.href = 'login.html';
    });
  }
  
  document.querySelectorAll('[data-filter]').forEach(button => {
    button.addEventListener('click', (event) => {
      event.preventDefault();
      currentFilter = event.currentTarget.dataset.filter;
      renderTasks();
    });
  });

  if (categoryFilterDropdown) {
    categoryFilterDropdown.addEventListener('click', (event) => {
      if (event.target.matches('[data-category-id]')) {
        event.preventDefault();
        currentCategoryFilter = event.target.dataset.categoryId;
        if(categoryFilterButton) {
            categoryFilterButton.textContent = event.target.textContent;
        }
        renderTasks();
      }
    });
  }
  
  if (editTaskDeleteButton) {
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
          editTaskErrorDiv.textContent = error.message;
          editTaskErrorDiv.classList.remove('d-none');
        }
      }
    });
  }
  
  const initializeApp = async () => {
    await fetchAllCategories();
    populateCategoryDropdowns();
    populateCategoryFilter();
    await fetchAllTasks();
  };

  initializeApp();
});