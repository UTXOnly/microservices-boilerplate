"""Pydantic schemas for Item resource."""

from datetime import datetime

from pydantic import BaseModel, ConfigDict, Field


class ItemBase(BaseModel):
    """Base schema with shared fields."""

    name: str = Field(..., min_length=1, max_length=255)
    description: str | None = None


class ItemCreate(ItemBase):
    """Schema for creating an item."""

    pass


class ItemUpdate(BaseModel):
    """Schema for partial updates (all fields optional)."""

    name: str | None = Field(None, min_length=1, max_length=255)
    description: str | None = None


class ItemRead(ItemBase):
    """Schema for reading an item from the API."""

    model_config = ConfigDict(from_attributes=True)

    id: int
    created_at: datetime
    updated_at: datetime
