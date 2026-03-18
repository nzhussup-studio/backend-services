package com.nzhussup.baseservice.service;

import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public abstract class BaseService<T> {

    protected final JpaRepository<T, Long> repository;

    public BaseService(JpaRepository<T, Long> repository) {
        this.repository = repository;
    }

    @Cacheable(value = "#root.targetClass.simpleName", key = "#root.targetClass.simpleName + '_all'")
    public List<T> findAll() {
        return repository.findAll();
    }

    public T findById(Long id) {
        return repository.findById(id).orElse(null);
    }

    @CacheEvict(value = "#root.targetClass.simpleName", key = "#root.targetClass.simpleName + '_all'")
    public T save(T entity) {
        return repository.save(entity);
    }

    @CacheEvict(value = "#root.targetClass.simpleName", key = "#root.targetClass.simpleName + '_all'")
    public void delete(Long id) {
        repository.deleteById(id);
    }

    @CacheEvict(value = "#root.targetClass.simpleName", key = "#root.targetClass.simpleName + '_all'")
    public T update(T entity) {
        return repository.save(entity);
    }
}
