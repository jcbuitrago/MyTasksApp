document.addEventListener('DOMContentLoaded', () => {

    const BASE_URL = 'http://localhost:8080';
    const authToken = localStorage.getItem('authToken');
  
    const categoryListContainer = document.getElementById('categoryListContainer');
    const taskListContainer = document.getElementById('taskListContainer');
    const taskListTitle = document.getElementById('task-list-title');
    const addCategoryForm = document.getElementById('addCategoryForm');
    const editCategoryForm = document.getElementById('editCategoryForm');
    const deleteCategoryButton = document.getElementById('deleteCategoryButton');
    const addModal = new bootstrap.Modal(document.getElementById('addCategoryModal'));
    const editModal = new bootstrap.Modal(document.getElementById('editCategoryModal'));
    
    let allCategories = [];
    let allTasks = [];
    let selectedCategoryId = null;

    const fetchAllTasks = async () => {
        try {
            const response = await fetch(`${BASE_URL}/tareas/usuario`, {
                headers: { 'Authorization': `Bearer ${authToken}` }
            });
            if (!response.ok) throw new Error('No se pudieron cargar las tareas.');
            allTasks = await response.json();
        } catch (error) {
            console.error(error);
        }
    };
  
    const fetchAndRenderCategories = async () => {
      try {
        const response = await fetch(`${BASE_URL}/categorias`, { headers: { 'Authorization': `Bearer ${authToken}` } });
        if (!response.ok) throw new Error('No se pudieron cargar las categorías.');
        allCategories = await response.json();
        
        // Limpiar solo las categorías cargadas dinámicamente
        document.querySelectorAll('#categoryListContainer .dynamic-category').forEach(el => el.remove());
        
        if (allCategories.length === 0) {
            // No hacer nada si no hay categorías, ya que "Sin Categoría" siempre está
        }
  
        allCategories.forEach(category => {
          const categoryElement = document.createElement('li');
          categoryElement.className = 'list-group-item dynamic-category d-flex align-items-center justify-content-between p-3 bg-transparent border-0 text-white';
          categoryElement.setAttribute('data-category-id', category.id);
          
          categoryElement.innerHTML = `
            <span class="m-0 text-truncate">${category.name}</span>
            <a href="#" class="edit-btn d-flex align-items-center justify-content-center text-decoration-none text-white-50" data-category-id="${category.id}" data-bs-toggle="modal" data-bs-target="#editCategoryModal">
              <svg xmlns="http://www.w3.org/2000/svg" width="16px" height="16px" fill="currentColor" viewBox="0 0 256 256"><path d="M227.31,73.37,182.63,28.68a16,16,0,0,0-22.63,0L36.69,152A15.86,15.86,0,0,0,32,163.31V208a16,16,0,0,0,16,16H92.69A15.86,15.86,0,0,0,104,219.31L227.31,96a16,16,0,0,0,0-22.63ZM92.69,208H48V163.31l88-88L180.69,120ZM192,108.68,147.31,64l24-24L216,84.68Z"></path></svg>
            </a>
          `;
          categoryListContainer.appendChild(categoryElement);
        });
      } catch (error) {
        console.error(error);
        categoryListContainer.innerHTML += `<p class="p-4 text-danger">${error.message}</p>`;
      }
    };

    const renderTasks = () => {
        if (selectedCategoryId === null) {
            taskListContainer.innerHTML = '<p class="text-center text-secondary p-5">Selecciona una categoría de la izquierda para ver sus tareas.</p>';
            taskListTitle.textContent = 'Tareas';
            return;
        }

        let filteredTasks;
        if (selectedCategoryId === 'unassigned') {
            filteredTasks = allTasks.filter(task => !task.category_id);
            taskListTitle.textContent = `Tareas "Sin Categoría"`;
        } else {
            filteredTasks = allTasks.filter(task => task.category_id == selectedCategoryId);
            const selectedCategory = allCategories.find(c => c.id == selectedCategoryId);
            taskListTitle.textContent = `Tareas en "${selectedCategory.name}"`;
        }
        
        taskListContainer.innerHTML = '';
        if (filteredTasks.length === 0) {
            taskListContainer.innerHTML = '<p class="text-center text-secondary p-5">No hay tareas en esta categoría.</p>';
            return;
        }

        filteredTasks.forEach(task => {
            const taskElement = document.createElement('div');
            taskElement.className = 'list-group-item bg-custom-dark text-white p-3 border-secondary-subtle';
            taskElement.innerHTML = `
                <h6 class="mb-1">${task.description}</h6>
                <small class="text-white-50">Vence: ${task.tentative_due_date} - Estado: ${task.status}</small>
            `;
            taskListContainer.appendChild(taskElement);
        });
    };

    categoryListContainer.addEventListener('click', (event) => {
        const categoryItem = event.target.closest('.list-group-item');
        const editButton = event.target.closest('.edit-btn');
        
        if (editButton) {
            const categoryId = editButton.dataset.categoryId;
            const category = allCategories.find(c => c.id == categoryId);
            if (category) {
                document.getElementById('editCategoryId').value = category.id;
                document.getElementById('editCategoryName').value = category.name;
                document.getElementById('editCategoryDescription').value = category.description;
            }
            return;
        }

        if (categoryItem) {
            selectedCategoryId = categoryItem.dataset.categoryId;
            document.querySelectorAll('#categoryListContainer .list-group-item').forEach(item => item.classList.remove('active'));
            categoryItem.classList.add('active');
            renderTasks();
        }
    });

    addCategoryForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        const name = document.getElementById('categoryName').value.trim();
        const description = document.getElementById('categoryDescription').value.trim();
        const errorDiv = document.getElementById('addCategoryError');
        if (!name) {
            errorDiv.textContent = 'El nombre es obligatorio.';
            errorDiv.classList.remove('d-none');
            return;
        }
        try {
            const response = await fetch(`${BASE_URL}/categorias`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${authToken}` },
                body: JSON.stringify({ name, description })
            });
            if (!response.ok) throw new Error('No se pudo crear la categoría.');
            addModal.hide();
            addCategoryForm.reset();
            await fetchAndRenderCategories();
        } catch (error) {
            errorDiv.textContent = error.message;
            errorDiv.classList.remove('d-none');
        }
    });

    editCategoryForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        const id = document.getElementById('editCategoryId').value;
        const name = document.getElementById('editCategoryName').value.trim();
        const description = document.getElementById('editCategoryDescription').value.trim();
        const errorDiv = document.getElementById('editCategoryError');
        try {
            const response = await fetch(`${BASE_URL}/categorias/${id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${authToken}` },
                body: JSON.stringify({ name, description })
            });
            if (!response.ok) throw new Error('No se pudo actualizar la categoría.');
            editModal.hide();
            await fetchAndRenderCategories();
        } catch (error) {
            errorDiv.textContent = error.message;
            errorDiv.classList.remove('d-none');
        }
    });

    deleteCategoryButton.addEventListener('click', async () => {
        const id = document.getElementById('editCategoryId').value;
        if (confirm('¿Estás seguro de que quieres eliminar esta categoría?')) {
            try {
                const response = await fetch(`${BASE_URL}/categorias/${id}`, {
                    method: 'DELETE',
                    headers: { 'Authorization': `Bearer ${authToken}` }
                });
                if (!response.ok) throw new Error('No se pudo eliminar la categoría.');
                editModal.hide();
                await fetchAndRenderCategories();
            } catch (error) {
                const errorDiv = document.getElementById('editCategoryError');
                errorDiv.textContent = error.message;
                errorDiv.classList.remove('d-none');
            }
        }
    });
    
    const initializeApp = async () => {
        await fetchAndRenderCategories();
        await fetchAllTasks();
    };

    initializeApp();
});