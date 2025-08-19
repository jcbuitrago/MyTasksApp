document.addEventListener('DOMContentLoaded', () => {

    const BASE_URL = 'http://localhost:8080';
    const authToken = localStorage.getItem('authToken');
  
    const categoryListContainer = document.getElementById('categoryListContainer');
    const addCategoryForm = document.getElementById('addCategoryForm');
    const editCategoryForm = document.getElementById('editCategoryForm');
    const deleteCategoryButton = document.getElementById('deleteCategoryButton');
    
    const addModal = new bootstrap.Modal(document.getElementById('addCategoryModal'));
    const editModal = new bootstrap.Modal(document.getElementById('editCategoryModal'));
    
    let allCategories = [];
  
    const fetchAndRenderCategories = async () => {
      try {
        const response = await fetch(`${BASE_URL}/categorias`, {
          headers: { 'Authorization': `Bearer ${authToken}` }
        });
        if (!response.ok) throw new Error('No se pudieron cargar las categorías.');
        
        allCategories = await response.json();
        
        categoryListContainer.innerHTML = ''; // Limpiar lista
        if (allCategories.length === 0) {
          categoryListContainer.innerHTML = '<p class="p-4 text-secondary">No hay categorías creadas.</p>';
          return;
        }
  
        allCategories.forEach(category => {
          const categoryElement = document.createElement('li');
          categoryElement.className = 'list-group-item d-flex align-items-center justify-content-between p-4 min-h-14 bg-custom-dark border-0';
          categoryElement.innerHTML = `
            <p class="m-0 text-truncate">${category.name}</p>
            <a href="#" class="edit-btn d-flex size-7 align-items-center justify-content-center text-decoration-none text-white" data-category-id="${category.id}" data-bs-toggle="modal" data-bs-target="#editCategoryModal">
              <svg xmlns="http://www.w3.org/2000/svg" width="24px" height="24px" fill="currentColor" viewBox="0 0 256 256"><path d="M227.31,73.37,182.63,28.68a16,16,0,0,0-22.63,0L36.69,152A15.86,15.86,0,0,0,32,163.31V208a16,16,0,0,0,16,16H92.69A15.86,15.86,0,0,0,104,219.31L227.31,96a16,16,0,0,0,0-22.63ZM92.69,208H48V163.31l88-88L180.69,120ZM192,108.68,147.31,64l24-24L216,84.68Z"></path></svg>
            </a>
          `;
          categoryListContainer.appendChild(categoryElement);
        });
      } catch (error) {
        console.error(error);
        categoryListContainer.innerHTML = `<p class="p-4 text-danger">${error.message}</p>`;
      }
    };
  
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
        await fetchAndRenderCategories(); // Recargar la lista
  
      } catch (error) {
        errorDiv.textContent = error.message;
        errorDiv.classList.remove('d-none');
      }
    });
  
    categoryListContainer.addEventListener('click', (event) => {
      const editButton = event.target.closest('.edit-btn');
      if (editButton) {
        const categoryId = editButton.dataset.categoryId;
        const category = allCategories.find(c => c.id == categoryId);
        
        if (category) {
          document.getElementById('editCategoryId').value = category.id;
          document.getElementById('editCategoryName').value = category.name;
          document.getElementById('editCategoryDescription').value = category.description;
        }
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
  
    fetchAndRenderCategories();
  });