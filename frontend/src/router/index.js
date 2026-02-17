import { createRouter, createWebHistory } from 'vue-router';
import Home from '@/views/Home.vue';
import Operations from '@/views/Operations.vue';
import Observability from '@/views/Observability.vue';
import Docs from '@/views/Docs.vue';

const defaultTitle = 'Enterprise Microservice Portal';

const routes = [
  {
    path: '/',
    name: 'home',
    component: Home,
    meta: {
      title: 'Enterprise Microservice Portal — Go + Vue Reference Platform',
      description: 'A production-focused reference platform demonstrating clean service boundaries, secure API patterns, resilient communication, and observability-first operations with Go and Vue.'
    }
  },
  {
    path: '/operations',
    name: 'operations',
    component: Operations,
    meta: {
      title: 'Operations — Enterprise Microservice Portal',
      description: 'Manage users, orders, and database migrations through the enterprise microservice operations dashboard.'
    }
  },
  {
    path: '/observability',
    name: 'observability',
    component: Observability,
    meta: {
      title: 'Observability — Enterprise Microservice Portal',
      description: 'Monitor service health endpoints, Prometheus metrics, and structured logging across enterprise microservices.'
    }
  },
  {
    path: '/docs',
    name: 'docs',
    component: Docs,
    meta: {
      title: 'API Docs — Enterprise Microservice Portal',
      description: 'Swagger documentation, Postman collections, and JWT security details for the enterprise microservice APIs.'
    }
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.afterEach((to) => {
  document.title = to.meta.title || defaultTitle;
  const descriptionTag = document.querySelector('meta[name="description"]');
  if (descriptionTag && to.meta.description) {
    descriptionTag.setAttribute('content', to.meta.description);
  }
});

export default router;
