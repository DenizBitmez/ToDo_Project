// API endpoint'ini dinamik olarak belirle
const API_BASE_URL = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
    ? 'http://localhost:8080'
    : window.location.origin;

// API istekleri için yardımcı fonksiyonlar
const api = {
    getToken() {
        return localStorage.getItem('token');
    },

    async get(url) {
        const token = this.getToken();
        if (!token) {
            throw new Error('Token bulunamadı. Lütfen tekrar giriş yapın.');
        }

        const response = await fetch(`${API_BASE_URL}${url}`, {
            headers: {
                'Authorization': `Bearer ${token}`,
                'Accept': 'application/json'
            },
            credentials: 'include'
        });

        if (response.status === 401) {
            localStorage.removeItem('token');
            localStorage.removeItem('role');
            localStorage.removeItem('username');
            window.location.href = '/login';
            throw new Error('Oturum süresi doldu. Lütfen tekrar giriş yapın.');
        }

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'İşlem başarısız');
        }

        return response.json();
    },

    async post(url, data) {
        const token = this.getToken();
        if (!token) {
            throw new Error('Token bulunamadı. Lütfen tekrar giriş yapın.');
        }

        const response = await fetch(`${API_BASE_URL}${url}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
                'Accept': 'application/json'
            },
            body: JSON.stringify(data),
            credentials: 'include'
        });

        if (response.status === 401) {
            localStorage.removeItem('token');
            localStorage.removeItem('role');
            localStorage.removeItem('username');
            window.location.href = '/login';
            throw new Error('Oturum süresi doldu. Lütfen tekrar giriş yapın.');
        }

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'İşlem başarısız');
        }

        return response.json();
    },

    async put(url, data) {
        const token = this.getToken();
        if (!token) {
            throw new Error('Token bulunamadı. Lütfen tekrar giriş yapın.');
        }

        const response = await fetch(`${API_BASE_URL}${url}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
                'Accept': 'application/json'
            },
            body: JSON.stringify(data),
            credentials: 'include'
        });

        if (response.status === 401) {
            localStorage.removeItem('token');
            localStorage.removeItem('role');
            localStorage.removeItem('username');
            window.location.href = '/login';
            throw new Error('Oturum süresi doldu. Lütfen tekrar giriş yapın.');
        }

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'İşlem başarısız');
        }

        return response.json();
    },

    async delete(url) {
        const token = this.getToken();
        if (!token) {
            throw new Error('Token bulunamadı. Lütfen tekrar giriş yapın.');
        }

        const response = await fetch(`${API_BASE_URL}${url}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Accept': 'application/json'
            },
            credentials: 'include'
        });

        if (response.status === 401) {
            localStorage.removeItem('token');
            localStorage.removeItem('role');
            localStorage.removeItem('username');
            window.location.href = '/login';
            throw new Error('Oturum süresi doldu. Lütfen tekrar giriş yapın.');
        }

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'İşlem başarısız');
        }

        return response.json();
    }
};

// Login işlemleri
document.addEventListener('DOMContentLoaded', () => {
    // Sayfa yüklendiğinde token kontrolü
    const token = localStorage.getItem('token');
    const role = localStorage.getItem('role');
    console.log('Sayfa yüklendi - Token:', token, 'Rol:', role);

    // Eğer token varsa ve doğru sayfada değilsek yönlendir
    if (token) {
        if (role === 'admin' && !window.location.href.includes('/admin')) {
            window.location.href = `${API_BASE_URL}/admin`;
        } else if (role === 'user' && !window.location.href.includes('/dashboard')) {
            window.location.href = `${API_BASE_URL}/dashboard`;
        }
    }

    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;

            try {
                console.log('Giriş denemesi:', username);
                const response = await fetch(`${API_BASE_URL}/token`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Accept': 'application/json'
                    },
                    body: JSON.stringify({ username, password }),
                    credentials: 'include'
                });

                console.log('Response status:', response.status);
                const contentType = response.headers.get('content-type');
                console.log('Content-Type:', contentType);

                if (!contentType || !contentType.includes('application/json')) {
                    throw new Error('Sunucu yanıtı JSON formatında değil');
                }

                const data = await response.json();
                console.log('Sunucu yanıtı:', data);

                if (!response.ok) {
                    throw new Error(data.error || 'Giriş başarısız');
                }

                if (data.token) {
                    console.log('Token kaydediliyor:', data.token);
                    localStorage.setItem('token', data.token);
                    localStorage.setItem('role', data.role);
                    localStorage.setItem('username', data.username);
                    console.log('Token kaydedildi, localStorage:', {
                        token: localStorage.getItem('token'),
                        role: localStorage.getItem('role'),
                        username: localStorage.getItem('username')
                    });

                    if (data.role === 'admin') {
                        console.log('Admin paneline yönlendiriliyor...');
                        window.location.href = `${API_BASE_URL}/admin`;
                    } else {
                        console.log('Dashboard\'a yönlendiriliyor...');
                        window.location.href = `${API_BASE_URL}/dashboard`;
                    }
                } else {
                    throw new Error('Token alınamadı');
                }
            } catch (error) {
                console.error('Giriş hatası:', error);
                alert('Giriş başarısız: ' + error.message);
            }
        });
    }
});

const todoOperations = {
    async loadTodos() {
        try {
            const token = localStorage.getItem('token');
            const role = localStorage.getItem('role');
            console.log('Todo\'lar yükleniyor... Token:', token, 'Rol:', role);
            
            if (!token) {
                throw new Error('Token bulunamadı. Lütfen tekrar giriş yapın.');
            }

            const response = await fetch('http://localhost:8080/todos', {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Accept': 'application/json'
                },
                credentials: 'include'
            });

            console.log('Response status:', response.status);
            
            if (response.status === 401) {
                localStorage.removeItem('token');
                localStorage.removeItem('role');
                localStorage.removeItem('username');
                window.location.href = '/login';
                throw new Error('Oturum süresi doldu. Lütfen tekrar giriş yapın.');
            }

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Todo\'lar yüklenemedi');
            }

            const todos = await response.json();
            console.log('Yüklenen todo\'lar:', todos);
            
            const todoList = document.getElementById('todoList');
            if (todoList) {
                if (!todos || !Array.isArray(todos)) {
                    console.log('Todo listesi boş veya geçersiz format');
                    todoList.innerHTML = '<div class="alert alert-info">Henüz todo eklenmemiş.</div>';
                    return;
                }

                if (todos.length === 0) {
                    console.log('Todo listesi boş');
                    todoList.innerHTML = '<div class="alert alert-info">Henüz todo eklenmemiş.</div>';
                    return;
                }

                todoList.innerHTML = todos.map(todo => `
                    <div class="todo-item" data-id="${todo.id}">
                        <div class="todo-header">
                            <div class="todo-title">
                                <h3>${todo.title}</h3>
                                <p class="text-muted">Kullanıcı: ${todo.username}</p>
                            </div>
                            <div class="todo-actions">
                                <button onclick="todoOperations.editTodo(${todo.id})" class="btn btn-primary">Düzenle</button>
                                <button onclick="todoOperations.deleteTodo(${todo.id})" class="btn btn-danger">Sil</button>
                            </div>
                        </div>
                        <div class="todo-steps">
                            <div class="steps-header">
                                <h4>Adımlar</h4>
                                <button onclick="todoStepOperations.showAddStepForm(${todo.id})" class="btn btn-success">Adım Ekle</button>
                            </div>
                            <div id="stepList-${todo.id}" class="step-list"></div>
                        </div>
                    </div>
                `).join('');

                todos.forEach(todo => {
                    todoStepOperations.loadSteps(todo.id);
                });
            }
        } catch (error) {
            console.error('Todo\'lar yüklenirken hata:', error);
            const todoList = document.getElementById('todoList');
            if (todoList) {
                todoList.innerHTML = `<div class="alert alert-danger">Todo'lar yüklenirken bir hata oluştu: ${error.message}</div>`;
            }
        }
    },

    async createTodo(title) {
        try {
            const token = localStorage.getItem('token');
            const role = localStorage.getItem('role');
            console.log('Yeni todo oluşturuluyor:', title, 'Token:', token, 'Rol:', role);
            
            if (!token) {
                throw new Error('Token bulunamadı. Lütfen tekrar giriş yapın.');
            }

            const response = await fetch('http://localhost:8080/todos', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                    'Accept': 'application/json'
                },
                body: JSON.stringify({ title }),
                credentials: 'include'
            });

            console.log('Response status:', response.status);

            if (response.status === 401) {
                localStorage.removeItem('token');
                localStorage.removeItem('role');
                localStorage.removeItem('username');
                window.location.href = '/login';
                throw new Error('Oturum süresi doldu. Lütfen tekrar giriş yapın.');
            }

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Todo oluşturulamadı');
            }

            const data = await response.json();
            console.log('Oluşturulan todo:', data);
            await this.loadTodos();
        } catch (error) {
            console.error('Todo oluşturulurken hata:', error);
            alert('Todo oluşturulurken bir hata oluştu: ' + error.message);
        }
    },

    async editTodo(id) {
        const newTitle = prompt('Yeni başlık:');
        if (newTitle) {
            try {
                const token = localStorage.getItem('token');
                console.log('Todo düzenleniyor:', id, newTitle, 'Token:', token);
                
                if (!token) {
                    throw new Error('Token bulunamadı. Lütfen tekrar giriş yapın.');
                }

                const response = await fetch(`http://localhost:8080/todos/${id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`,
                        'Accept': 'application/json'
                    },
                    body: JSON.stringify({ title: newTitle }),
                    credentials: 'include'
                });

                console.log('Response status:', response.status);

                if (response.status === 401) {
                    localStorage.removeItem('token');
                    localStorage.removeItem('role');
                    localStorage.removeItem('username');
                    window.location.href = '/login';
                    throw new Error('Oturum süresi doldu. Lütfen tekrar giriş yapın.');
                }

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || 'Todo düzenlenemedi');
                }

                const data = await response.json();
                console.log('Düzenlenen todo:', data);
                await this.loadTodos();
            } catch (error) {
                console.error('Todo düzenlenirken hata:', error);
                alert('Todo düzenlenirken bir hata oluştu: ' + error.message);
            }
        }
    },

    async deleteTodo(id) {
        if (confirm('Bu todo\'yu silmek istediğinizden emin misiniz?')) {
            try {
                const token = localStorage.getItem('token');
                console.log('Todo siliniyor:', id, 'Token:', token);
                
                if (!token) {
                    throw new Error('Token bulunamadı. Lütfen tekrar giriş yapın.');
                }

                const response = await fetch(`http://localhost:8080/todos/${id}`, {
                    method: 'DELETE',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Accept': 'application/json'
                    },
                    credentials: 'include'
                });

                console.log('Response status:', response.status);

                if (response.status === 401) {
                    localStorage.removeItem('token');
                    localStorage.removeItem('role');
                    localStorage.removeItem('username');
                    window.location.href = '/login';
                    throw new Error('Oturum süresi doldu. Lütfen tekrar giriş yapın.');
                }

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || 'Todo silinemedi');
                }

                const data = await response.json();
                console.log('Silinen todo:', data);
                await this.loadTodos();
            } catch (error) {
                console.error('Todo silinirken hata:', error);
                alert('Todo silinirken bir hata oluştu: ' + error.message);
            }
        }
    }
};

const todoStepOperations = {
    async loadSteps(todoId) {
        try {
            const token = localStorage.getItem('token');
            console.log('Todo adımları yükleniyor... Todo ID:', todoId, 'Token:', token);
            
            if (!token) {
                throw new Error('Token bulunamadı. Lütfen tekrar giriş yapın.');
            }

            const response = await fetch(`http://localhost:8080/todos/${todoId}/steps`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Accept': 'application/json'
                },
                credentials: 'include'
            });

            console.log('Response status:', response.status);
            console.log('Response headers:', Object.fromEntries(response.headers.entries()));
            
            // Ham yanıtı kontrol et
            const rawResponse = await response.text();
            console.log('Raw response:', rawResponse);
            
            if (response.status === 401) {
                localStorage.removeItem('token');
                localStorage.removeItem('role');
                localStorage.removeItem('username');
                window.location.href = '/login';
                throw new Error('Oturum süresi doldu. Lütfen tekrar giriş yapın.');
            }

            if (!response.ok) {
                throw new Error(`Adımlar yüklenemedi: ${response.status} ${response.statusText}`);
            }

            // JSON parse etmeyi dene
            let steps;
            try {
                steps = JSON.parse(rawResponse);
            } catch (parseError) {
                console.error('JSON parse hatası:', parseError);
                throw new Error('Sunucu yanıtı geçerli bir JSON formatında değil');
            }

            console.log('Yüklenen adımlar:', steps);
            
            const stepList = document.getElementById(`stepList-${todoId}`);
            if (stepList) {
                if (!steps || !Array.isArray(steps)) {
                    console.log('Adım listesi boş veya geçersiz format');
                    stepList.innerHTML = '<div class="alert alert-info">Henüz adım eklenmemiş.</div>';
                    return;
                }

                if (steps.length === 0) {
                    console.log('Adım listesi boş');
                    stepList.innerHTML = '<div class="alert alert-info">Henüz adım eklenmemiş.</div>';
                    return;
                }

                stepList.innerHTML = steps.map(step => `
                    <div class="step-item" data-id="${step.id}">
                        <div class="step-content">
                            <span class="step-title">${step.content || step.title}</span>
                            <span class="step-status ${step.status === 1 ? 'completed' : ''}">
                                ${step.status === 1 ? 'Tamamlandı' : 'Devam Ediyor'}
                            </span>
                        </div>
                        <div class="step-actions">
                            <button onclick="todoStepOperations.toggleStep(${todoId}, ${step.id})" 
                                    class="btn btn-${step.status === 1 ? 'warning' : 'success'}">
                                ${step.status === 1 ? 'Geri Al' : 'Tamamla'}
                            </button>
                            <button onclick="todoStepOperations.editStep(${todoId}, ${step.id})" 
                                    class="btn btn-primary">Düzenle</button>
                            <button onclick="todoStepOperations.deleteStep(${todoId}, ${step.id})" 
                                    class="btn btn-danger">Sil</button>
                        </div>
                    </div>
                `).join('');
            }
        } catch (error) {
            console.error('Adımlar yüklenirken hata:', error);
            const stepList = document.getElementById(`stepList-${todoId}`);
            if (stepList) {
                stepList.innerHTML = `<div class="alert alert-danger">Adımlar yüklenirken bir hata oluştu: ${error.message}</div>`;
            }
        }
    },

    async createStep(todoId, title) {
        try {
            const token = localStorage.getItem('token');
            console.log('Yeni adım oluşturuluyor:', title, 'Todo ID:', todoId, 'Token:', token);
            
            if (!token) {
                throw new Error('Token bulunamadı. Lütfen tekrar giriş yapın.');
            }

            const response = await fetch(`http://localhost:8080/todos/${todoId}/steps`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                    'Accept': 'application/json'
                },
                body: JSON.stringify({ title }),
                credentials: 'include'
            });

            console.log('Response status:', response.status);

            if (response.status === 401) {
                localStorage.removeItem('token');
                localStorage.removeItem('role');
                localStorage.removeItem('username');
                window.location.href = '/login';
                throw new Error('Oturum süresi doldu. Lütfen tekrar giriş yapın.');
            }

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Adım oluşturulamadı');
            }

            const data = await response.json();
            console.log('Oluşturulan adım:', data);
            await this.loadSteps(todoId);
        } catch (error) {
            console.error('Adım oluşturulurken hata:', error);
            alert('Adım oluşturulurken bir hata oluştu: ' + error.message);
        }
    },

    async editStep(todoId, stepId) {
        const newTitle = prompt('Yeni adım başlığı:');
        if (newTitle) {
            try {
                const token = localStorage.getItem('token');
                console.log('Adım düzenleniyor:', stepId, newTitle, 'Todo ID:', todoId, 'Token:', token);
                
                if (!token) {
                    throw new Error('Token bulunamadı. Lütfen tekrar giriş yapın.');
                }

                const response = await fetch(`http://localhost:8080/todos/${todoId}/steps/${stepId}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`,
                        'Accept': 'application/json'
                    },
                    body: JSON.stringify({ title: newTitle }),
                    credentials: 'include'
                });

                console.log('Response status:', response.status);

                if (response.status === 401) {
                    localStorage.removeItem('token');
                    localStorage.removeItem('role');
                    localStorage.removeItem('username');
                    window.location.href = '/login';
                    throw new Error('Oturum süresi doldu. Lütfen tekrar giriş yapın.');
                }

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || 'Adım düzenlenemedi');
                }

                const data = await response.json();
                console.log('Düzenlenen adım:', data);
                await this.loadSteps(todoId);
            } catch (error) {
                console.error('Adım düzenlenirken hata:', error);
                alert('Adım düzenlenirken bir hata oluştu: ' + error.message);
            }
        }
    },

    async deleteStep(todoId, stepId) {
        if (confirm('Bu adımı silmek istediğinizden emin misiniz?')) {
            try {
                const token = localStorage.getItem('token');
                console.log('Adım siliniyor:', stepId, 'Todo ID:', todoId, 'Token:', token);
                
                if (!token) {
                    throw new Error('Token bulunamadı. Lütfen tekrar giriş yapın.');
                }

                const response = await fetch(`http://localhost:8080/todos/${todoId}/steps/${stepId}`, {
                    method: 'DELETE',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Accept': 'application/json'
                    },
                    credentials: 'include'
                });

                console.log('Response status:', response.status);

                if (response.status === 401) {
                    localStorage.removeItem('token');
                    localStorage.removeItem('role');
                    localStorage.removeItem('username');
                    window.location.href = '/login';
                    throw new Error('Oturum süresi doldu. Lütfen tekrar giriş yapın.');
                }

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || 'Adım silinemedi');
                }

                const data = await response.json();
                console.log('Silinen adım:', data);
                await this.loadSteps(todoId);
            } catch (error) {
                console.error('Adım silinirken hata:', error);
                alert('Adım silinirken bir hata oluştu: ' + error.message);
            }
        }
    },

    async toggleStep(todoId, stepId) {
        try {
            const token = localStorage.getItem('token');
            console.log('Adım durumu değiştiriliyor:', stepId, 'Todo ID:', todoId, 'Token:', token);
            
            if (!token) {
                throw new Error('Token bulunamadı. Lütfen tekrar giriş yapın.');
            }

            const response = await fetch(`http://localhost:8080/todos/${todoId}/steps/${stepId}/toggle`, {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                },
                credentials: 'include'
            });

            console.log('Response status:', response.status);
            console.log('Response headers:', Object.fromEntries(response.headers.entries()));

            if (response.status === 401) {
                localStorage.removeItem('token');
                localStorage.removeItem('role');
                localStorage.removeItem('username');
                window.location.href = '/login';
                throw new Error('Oturum süresi doldu. Lütfen tekrar giriş yapın.');
            }

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Adım durumu değiştirilemedi');
            }

            const data = await response.json();
            console.log('Durumu değiştirilen adım:', data);

            await this.loadSteps(todoId);
        } catch (error) {
            console.error('Adım durumu değiştirilirken hata:', error);
            alert('Adım durumu değiştirilirken bir hata oluştu: ' + error.message);
        }
    },

    showAddStepForm(todoId) {
        const stepList = document.getElementById(`stepList-${todoId}`);
        if (stepList) {
            const formHtml = `
                <div class="add-step-form">
                    <input type="text" id="newStepTitle-${todoId}" class="form-control" placeholder="Adım başlığı">
                    <button onclick="todoStepOperations.submitNewStep(${todoId})" class="btn btn-primary">Ekle</button>
                    <button onclick="todoStepOperations.cancelAddStep(${todoId})" class="btn btn-secondary">İptal</button>
                </div>
            `;
            stepList.insertAdjacentHTML('afterbegin', formHtml);
        }
    },

    async submitNewStep(todoId) {
        const titleInput = document.getElementById(`newStepTitle-${todoId}`);
        const title = titleInput.value.trim();
        
        if (title) {
            await this.createStep(todoId, title);
            this.cancelAddStep(todoId);
        } else {
            alert('Lütfen bir adım başlığı girin');
        }
    },

    cancelAddStep(todoId) {
        const form = document.querySelector(`#stepList-${todoId} .add-step-form`);
        if (form) {
            form.remove();
        }
    }
};

document.addEventListener('DOMContentLoaded', () => {
    const todoList = document.getElementById('todoList');
    if (todoList) {
        console.log('Todo listesi bulundu, yükleniyor...');
        todoOperations.loadTodos();
    } else {
        console.log('Todo listesi bulunamadı');
    }
});

const createTodoForm = document.getElementById('createTodoForm');
if (createTodoForm) {
    createTodoForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const titleInput = document.getElementById('newTodoTitle');
        const title = titleInput.value.trim();
        
        if (title) {
            console.log('Form gönderiliyor, başlık:', title);
            console.log('Mevcut token:', localStorage.getItem('token'));
            console.log('Mevcut rol:', localStorage.getItem('role'));
            await todoOperations.createTodo(title);
            titleInput.value = '';
        } else {
            alert('Lütfen bir başlık girin');
        }
    });
}

const logoutButton = document.getElementById('logoutButton');
if (logoutButton) {
    logoutButton.addEventListener('click', () => {
        console.log('Çıkış yapılıyor...');
        localStorage.removeItem('token');
        localStorage.removeItem('role');
        localStorage.removeItem('username');
        window.location.href = '/login';
    });
}