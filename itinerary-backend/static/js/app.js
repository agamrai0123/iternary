// ==================== API Utility Functions ====================

async function apiCall(method, endpoint, data = null) {
    const options = {
        method: method,
        headers: {
            'Content-Type': 'application/json',
        }
    };

    if (data) {
        options.body = JSON.stringify(data);
    }

    try {
        const response = await fetch(endpoint, options);
        if (!response.ok) {
            throw new Error(`API error: ${response.statusText}`);
        }
        return await response.json();
    } catch (error) {
        console.error('API call failed:', error);
        showNotification('Error: ' + error.message, 'error');
        throw error;
    }
}

// ==================== Like Functionality ====================

async function likeItinerary(itineraryId) {
    try {
        const result = await apiCall('POST', `/api/itineraries/${itineraryId}/like`);
        showNotification('Thanks for liking this itinerary! ❤️', 'success');
        
        // Reload page to show updated like count
        setTimeout(() => location.reload(), 500);
    } catch (error) {
        showNotification('Failed to like itinerary', 'error');
    }
}

// ==================== Comment Functionality ====================

async function postComment(event, itineraryId) {
    event.preventDefault();
    
    const form = event.target;
    const content = form.querySelector('textarea[name="content"]').value;
    
    if (!content.trim()) {
        showNotification('Comment cannot be empty', 'error');
        return;
    }

    try {
        const result = await apiCall('POST', `/api/itineraries/${itineraryId}/comments`, {
            content: content,
            user_id: 'user-001' // TODO: Get from session/auth
        });
        
        showNotification('Comment posted successfully! 🎉', 'success');
        form.reset();
        
        // Reload comments
        // loadComments(itineraryId);
    } catch (error) {
        showNotification('Failed to post comment', 'error');
    }
}

// ==================== Copy to Plans ====================

function copyPlan(itineraryId) {
    showNotification('Feature coming soon! You\'ll be able to copy this plan to your profile.', 'info');
    // TODO: Implement plan copying
}

// ==================== Notifications ====================

function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.textContent = message;
    
    document.body.appendChild(notification);
    
    // Animation
    setTimeout(() => notification.classList.add('show'), 10);
    
    // Auto remove
    setTimeout(() => {
        notification.classList.remove('show');
        setTimeout(() => notification.remove(), 300);
    }, 3000);
}

// ==================== Form Utilities ====================

function validateForm(formId) {
    const form = document.getElementById(formId);
    if (!form) return true;
    
    return form.checkValidity();
}

function formatPrice(price) {
    return new Intl.NumberFormat('en-IN', {
        style: 'currency',
        currency: 'INR',
        minimumFractionDigits: 0
    }).format(price);
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-IN', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });
}

// ==================== Template Helpers ====================

function typeIcon(type) {
    const icons = {
        'stay': '🏨',
        'food': '🍽️',
        'activity': '🎯',
        'transport': '🚗',
        'other': '📌'
    };
    return icons[type] || '📌';
}

function toUpper(str) {
    return str.toUpperCase();
}

// ==================== Search Utilities ====================

function performSearch() {
    const query = document.getElementById('q')?.value || '';
    const destination = document.getElementById('destination')?.value || '';
    const maxBudget = document.getElementById('max_budget')?.value || '';
    
    const params = new URLSearchParams();
    if (query) params.append('q', query);
    if (destination) params.append('destination', destination);
    if (maxBudget) params.append('max_budget', maxBudget);
    
    window.location.href = `/search?${params.toString()}`;
}

// ==================== DOM Ready ====================

document.addEventListener('DOMContentLoaded', function() {
    // Initialize tooltips if you use them
    initializeTooltips();
    
    // Add keyboard shortcuts
    setupKeyboardShortcuts();
});

function initializeTooltips() {
    // Add title attributes as tooltips
    const elements = document.querySelectorAll('[data-tooltip]');
    elements.forEach(el => {
        el.title = el.dataset.tooltip;
    });
}

function setupKeyboardShortcuts() {
    document.addEventListener('keydown', function(event) {
        // Ctrl/Cmd + K: Focus search
        if ((event.ctrlKey || event.metaKey) && event.key === 'k') {
            event.preventDefault();
            const searchInput = document.getElementById('q');
            if (searchInput) searchInput.focus();
        }
        
        // Ctrl/Cmd + H: Go home
        if ((event.ctrlKey || event.metaKey) && event.key === 'h') {
            event.preventDefault();
            window.location.href = '/';
        }
    });
}

// ==================== Notification Styles (CSS) ====================

// Add this to your CSS or create a separate file:
/*
.notification {
    position: fixed;
    top: 20px;
    right: 20px;
    padding: 1rem 1.5rem;
    border-radius: 4px;
    color: white;
    font-weight: 500;
    z-index: 1000;
    opacity: 0;
    transition: opacity 0.3s ease;
}

.notification.show {
    opacity: 1;
}

.notification-success {
    background-color: #27ae60;
}

.notification-error {
    background-color: #e74c3c;
}

.notification-info {
    background-color: #3498db;
}

@media (max-width: 768px) {
    .notification {
        left: 10px;
        right: 10px;
    }
}
*/

// ==================== Export for use in templates ====================

window.likeItinerary = likeItinerary;
window.postComment = postComment;
window.copyPlan = copyPlan;
window.typeIcon = typeIcon;
window.toUpper = toUpper;
window.performSearch = performSearch;
