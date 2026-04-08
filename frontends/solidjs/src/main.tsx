import { render } from 'solid-js/web';
import { AppRouter } from './router';
import '../../../shared/design-system/index.css';

const root = document.getElementById('root');

if (!root) {
  throw new Error('Root element not found');
}

render(() => <AppRouter />, root);
