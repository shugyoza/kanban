import { ApplicationConfig, provideBrowserGlobalErrorListeners } from '@angular/core';
import { provideRouter } from '@angular/router';
import { routes } from './app.routes';

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideRouter(routes)

    // ONLY include this if you are actively passing custom feature arguments!
    // provideHttpClient(withInterceptors([authInterceptor])) 
  ]
};
