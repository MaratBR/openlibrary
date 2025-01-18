import Alpine from 'alpinejs'
//@ts-expect-error
import ajax from '@imacrayon/alpine-ajax'

(window as any).Alpine = Alpine;
Alpine.plugin(ajax)
Alpine.start();
