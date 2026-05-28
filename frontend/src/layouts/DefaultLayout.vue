<script setup lang="ts">
import { RouterView, RouterLink } from 'vue-router'
</script>

<template>
  <div class="layout">

    <!-- HEADER -->
    <header class="header">
      <div class="logo">
        My App
      </div>

      <nav class="nav">
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/about">About</RouterLink>
        <RouterLink to="/login">Login</RouterLink>
      </nav>
    </header>

    <!-- CONTENT -->
    <main class="content">
      <RouterView />
    </main>

    <!-- FOOTER -->
    <footer class="footer">
      © 2026 My App
    </footer>

  </div>
</template>

<style scoped>
.layout {
  display: flex;
  flex-direction: column;
  width: 100%;
  min-height: 100vh;
  margin: 0 auto;
  box-sizing: border-box;
  overflow-x: hidden;
}

/* HEADER */
.header {
  flex: 0 0 auto;
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  background: #42b983;
  color: white;
  box-sizing: border-box;
  gap: 1rem;
}

.logo {
  white-space: nowrap;
  font-weight: bold;
  font-size: 1.25rem;
}

.nav {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
  justify-content: flex-end;
  flex: 1;
}

.nav a {
  color: white;
  text-decoration: none;
  padding: 0.5rem 0.75rem;
  border-radius: 4px;
  transition: background-color 0.3s;
  white-space: nowrap;
}

.nav a:hover {
  background-color: rgba(255, 255, 255, 0.2);
}

/* CONTENT */
.content {
  flex: 1 1 auto;
  width: 100%;
  padding: 1.5rem;
  background: #f5f5f5;
  box-sizing: border-box;
  overflow-y: auto;
}

/* FOOTER */
.footer {
  flex: 0 0 auto;
  text-align: center;
  padding: 0.75rem;
  background: #eee;
  box-sizing: border-box;
  width: 100%;
  border-top: 1px solid #ddd;
}

/* Responsive scaling */
@media (max-width: 768px) {
  .header {
    padding: 0.75rem 1rem;
    gap: 0.75rem;
  }
  
  .logo {
    font-size: 1.1rem;
  }
  
  .nav {
    gap: 0.5rem;
  }
  
  .content {
    padding: 1rem;
  }
}

@media (max-width: 480px) {
  .header {
    padding: 0.5rem 0.75rem;
    flex-wrap: wrap;
    gap: 0.5rem;
  }
  
  .logo {
    flex-basis: 100%;
    font-size: 1rem;
  }
  
  .nav {
    gap: 0.5rem;
    width: 100%;
    justify-content: flex-start;
  }
  
  .nav a {
    padding: 0.4rem 0.6rem;
    font-size: 0.9rem;
  }
  
  .content {
    padding: 0.75rem;
  }
  
  .footer {
    padding: 0.5rem;
    font-size: 0.85rem;
  }
}
</style>