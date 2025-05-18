# Используем официальный образ Nginx
FROM nginx:alpine

# Удаляем стандартный конфиг
RUN rm /etc/nginx/conf.d/default.conf

# Копируем наш конфигурационный файл
COPY nginx.conf /etc/nginx/nginx.conf

# Открываем порт 80
EXPOSE 7070

# Запускаем Nginx
CMD ["nginx", "-g", "daemon off;"]
