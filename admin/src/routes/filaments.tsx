import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import {
  useCreateFilament,
  useDeleteFilament,
  useFilaments,
  useUpdateFilament,
} from '../lib/queries'
import type { Filament } from '../lib/api'

export const Route = createFileRoute('/filaments')({ component: FilamentsPage })

interface FilamentForm {
  name: string
  amount_grams: string
  total_price: string
  color_hex: string
  filament_type: string
  temperature: string
  bed_temperature: string
  first_layer_temperature: string
  first_layer_bed_temperature: string
  filament_diameter: string
  extrusion_multiplier: string
  filament_density: string
}

const emptyForm: FilamentForm = {
  name: '',
  amount_grams: '',
  total_price: '',
  color_hex: '#ffffff',
  filament_type: 'PLA',
  temperature: '',
  bed_temperature: '',
  first_layer_temperature: '',
  first_layer_bed_temperature: '',
  filament_diameter: '',
  extrusion_multiplier: '',
  filament_density: '',
}

function costPerGram(f: FilamentForm): string {
  const a = parseFloat(f.amount_grams)
  const p = parseFloat(f.total_price)
  if (!a || !p || a <= 0) return '—'
  return (p / a).toFixed(4)
}

function FilamentsPage() {
  const { data: filaments, isLoading } = useFilaments()
  const create = useCreateFilament()
  const update = useUpdateFilament()
  const remove = useDeleteFilament()

  const [addForm, setAddForm] = useState<FilamentForm | null>(null)
  const [editId, setEditId] = useState<string | null>(null)
  const [editForm, setEditForm] = useState<FilamentForm>(emptyForm)

  function openAdd() {
    setEditId(null)
    setAddForm(emptyForm)
  }

  function openEdit(f: Filament) {
    setAddForm(null)
    setEditId(f.id)
    const totalPrice = (f.cost_per_gram * f.current_amount).toFixed(2)
    setEditForm({
      name: f.filament_name,
      amount_grams: String(f.current_amount),
      total_price: totalPrice,
      color_hex: f.color_hex.startsWith('#') ? f.color_hex : `#${f.color_hex}`,
      filament_type: f.filament_type,
      temperature: f.temperature,
      bed_temperature: f.bed_temperature,
      first_layer_temperature: f.first_layer_temperature,
      first_layer_bed_temperature: f.first_layer_bed_temperature,
      filament_diameter: f.filament_diameter,
      extrusion_multiplier: f.extrusion_multiplier,
      filament_density: f.filament_density,
    })
  }

  function handleAdd(e: React.FormEvent) {
    e.preventDefault()
    if (!addForm) return
    create.mutate(
      {
        name: addForm.name,
        amount_grams: parseInt(addForm.amount_grams),
        total_price: parseFloat(addForm.total_price),
        color_hex: addForm.color_hex.replace(/^#/, ''),
        filament_type: addForm.filament_type,
        temperature: addForm.temperature,
        bed_temperature: addForm.bed_temperature,
        first_layer_temperature: addForm.first_layer_temperature,
        first_layer_bed_temperature: addForm.first_layer_bed_temperature,
        filament_diameter: addForm.filament_diameter,
        extrusion_multiplier: addForm.extrusion_multiplier,
        filament_density: addForm.filament_density,
      },
      { onSuccess: () => setAddForm(null) },
    )
  }

  function handleEdit(e: React.FormEvent) {
    e.preventDefault()
    if (!editId) return
    update.mutate(
      {
        id: editId,
        name: editForm.name,
        amount_grams: parseInt(editForm.amount_grams),
        total_price: parseFloat(editForm.total_price),
        color_hex: editForm.color_hex.replace(/^#/, ''),
        filament_type: editForm.filament_type || undefined,
        temperature: editForm.temperature || undefined,
        bed_temperature: editForm.bed_temperature || undefined,
        first_layer_temperature: editForm.first_layer_temperature || undefined,
        first_layer_bed_temperature: editForm.first_layer_bed_temperature || undefined,
        filament_diameter: editForm.filament_diameter || undefined,
        extrusion_multiplier: editForm.extrusion_multiplier || undefined,
        filament_density: editForm.filament_density || undefined,
      },
      { onSuccess: () => setEditId(null) },
    )
  }

  function handleDelete(id: string) {
    if (!confirm('Delete this filament?')) return
    remove.mutate(id)
  }

  if (isLoading) return <p>Loading…</p>

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-2xl font-bold">Filaments</h1>
        {!addForm && (
          <button
            onClick={openAdd}
            className="hover-black border-2 border-black bg-white hover:cursor-pointer px-6 py-2 font-semibold"
          >
            Add Filament
          </button>
        )}
      </div>

      <table className="w-full border-2 border-black text-left">
        <thead className="border-b-2 border-black bg-gray-100">
          <tr>
            <th className="px-4 py-2 font-semibold">Color</th>
            <th className="px-4 py-2 font-semibold">Name</th>
            <th className="px-4 py-2 font-semibold">Stock (g)</th>
            <th className="px-4 py-2 font-semibold">Cost/g</th>
            <th className="px-4 py-2 font-semibold">Actions</th>
          </tr>
        </thead>
        <tbody>
          {filaments?.map((f) =>
            editId === f.id ? (
              <tr key={f.id} className="border-b border-black">
                <td colSpan={5} className="px-4 py-3">
                  <form onSubmit={handleEdit} className="flex flex-col gap-6">
                    <div className="flex flex-wrap items-end gap-3">
                      <div className="flex flex-col gap-1">
                        <label className="text-xs font-semibold">Name</label>
                        <input
                          value={editForm.name}
                          onChange={(e) =>
                            setEditForm({ ...editForm, name: e.target.value })
                          }
                          required
                          className="border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                        />
                      </div>
                      <div className="flex flex-col gap-1">
                        <label className="text-xs font-semibold">Stock (g)</label>
                        <input
                          type="number"
                          min="1"
                          value={editForm.amount_grams}
                          onChange={(e) =>
                            setEditForm({ ...editForm, amount_grams: e.target.value })
                          }
                          required
                          className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                        />
                      </div>
                      <div className="flex flex-col gap-1">
                        <label className="text-xs font-semibold">Spool Price (€)</label>
                        <input
                          type="number"
                          min="0.01"
                          step="0.01"
                          value={editForm.total_price}
                          onChange={(e) =>
                            setEditForm({ ...editForm, total_price: e.target.value })
                          }
                          required
                          className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                        />
                      </div>
                      <div className="flex flex-col gap-1">
                        <label className="text-xs font-semibold">Cost/g</label>
                        <span className="border-2 border-black bg-gray-100 px-3 py-1.5 text-sm">
                          €{costPerGram(editForm)}
                        </span>
                      </div>
                      <div className="flex flex-col gap-1">
                        <label className="text-xs font-semibold">Color</label>
                        <input
                          type="color"
                          value={editForm.color_hex}
                          onChange={(e) =>
                            setEditForm({ ...editForm, color_hex: e.target.value })
                          }
                          className="h-[38px] w-16 border-2 border-black bg-white p-0"
                        />
                      </div>
                      <div className="flex flex-col gap-1">
                        <label className="text-xs font-semibold">Filament Type</label>
                        <input
                          value={editForm.filament_type}
                          onChange={(e) =>
                            setEditForm({ ...editForm, filament_type: e.target.value })
                          }
                          placeholder="PLA"
                          className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                        />
                      </div>
                    </div>

                    <fieldset className="flex flex-col gap-3 border-2 border-black p-4">
                      <legend className="px-2 text-xs font-semibold">
                        Slicer profile (leave blank to keep existing values)
                      </legend>
                      <div className="flex flex-wrap gap-3">
                        <div className="flex flex-col gap-1">
                          <label className="text-xs font-semibold">Nozzle Temp (°C)</label>
                          <input
                            type="number"
                            value={editForm.temperature}
                            onChange={(e) =>
                              setEditForm({ ...editForm, temperature: e.target.value })
                            }
                            placeholder="220"
                            className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                          />
                        </div>
                        <div className="flex flex-col gap-1">
                          <label className="text-xs font-semibold">Bed Temp (°C)</label>
                          <input
                            type="number"
                            value={editForm.bed_temperature}
                            onChange={(e) =>
                              setEditForm({ ...editForm, bed_temperature: e.target.value })
                            }
                            placeholder="60"
                            className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                          />
                        </div>
                        <div className="flex flex-col gap-1">
                          <label className="text-xs font-semibold">
                            First Layer Nozzle (°C)
                          </label>
                          <input
                            type="number"
                            value={editForm.first_layer_temperature}
                            onChange={(e) =>
                              setEditForm({
                                ...editForm,
                                first_layer_temperature: e.target.value,
                              })
                            }
                            placeholder="225"
                            className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                          />
                        </div>
                        <div className="flex flex-col gap-1">
                          <label className="text-xs font-semibold">
                            First Layer Bed (°C)
                          </label>
                          <input
                            type="number"
                            value={editForm.first_layer_bed_temperature}
                            onChange={(e) =>
                              setEditForm({
                                ...editForm,
                                first_layer_bed_temperature: e.target.value,
                              })
                            }
                            placeholder="65"
                            className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                          />
                        </div>
                        <div className="flex flex-col gap-1">
                          <label className="text-xs font-semibold">Diameter (mm)</label>
                          <input
                            type="number"
                            step="0.01"
                            value={editForm.filament_diameter}
                            onChange={(e) =>
                              setEditForm({
                                ...editForm,
                                filament_diameter: e.target.value,
                              })
                            }
                            placeholder="1.75"
                            className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                          />
                        </div>
                        <div className="flex flex-col gap-1">
                          <label className="text-xs font-semibold">Extrusion Mult.</label>
                          <input
                            type="number"
                            step="0.01"
                            value={editForm.extrusion_multiplier}
                            onChange={(e) =>
                              setEditForm({
                                ...editForm,
                                extrusion_multiplier: e.target.value,
                              })
                            }
                            placeholder="1"
                            className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                          />
                        </div>
                        <div className="flex flex-col gap-1">
                          <label className="text-xs font-semibold">Density (g/cm³)</label>
                          <input
                            type="number"
                            step="0.01"
                            value={editForm.filament_density}
                            onChange={(e) =>
                              setEditForm({
                                ...editForm,
                                filament_density: e.target.value,
                              })
                            }
                            placeholder="1.24"
                            className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                          />
                        </div>
                      </div>
                    </fieldset>

                    <div className="flex gap-3">
                      <button
                        type="submit"
                        disabled={update.isPending}
                        className="hover-black border-2 border-black bg-white hover:cursor-pointer px-4 py-1.5 font-semibold disabled:opacity-50"
                      >
                        Save
                      </button>
                      <button
                        type="button"
                        onClick={() => setEditId(null)}
                        className="border-2 border-black bg-white hover:cursor-pointer px-4 py-1.5 font-semibold hover:bg-gray-100"
                      >
                        Cancel
                      </button>
                    </div>
                  </form>
                </td>
              </tr>
            ) : (
              <tr key={f.id} className="border-b border-black last:border-b-0">
                <td className="px-4 py-2">
                  <span
                    className="inline-block h-6 w-6 border-2 border-black align-middle"
                    style={{
                      backgroundColor: f.color_hex.startsWith('#')
                        ? f.color_hex
                        : `#${f.color_hex}`,
                    }}
                    title={f.color_hex}
                  />
                </td>
                <td className="px-4 py-2">{f.filament_name}</td>
                <td className="px-4 py-2">{f.current_amount}</td>
                <td className="px-4 py-2">€{f.cost_per_gram.toFixed(4)}</td>
                <td className="px-4 py-2">
                  <div className="flex gap-2">
                    <button
                      onClick={() => openEdit(f)}
                      className="hover-black border-2 border-black bg-white hover:cursor-pointer px-3 py-1 text-sm font-semibold"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDelete(f.id)}
                      className="hover-black border-2 border-black bg-red-200 hover:cursor-pointer hover:bg-red-400 px-3 py-1 text-sm font-semibold"
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            ),
          )}
          {filaments?.length === 0 && (
            <tr>
              <td colSpan={5} className="px-4 py-6 text-center text-gray-500">
                No filaments
              </td>
            </tr>
          )}
        </tbody>
      </table>

      {addForm !== null && (
        <div className="mt-6 border-2 border-black p-6">
          <h2 className="mb-4 font-bold">Add Filament</h2>
          <form onSubmit={handleAdd} className="flex flex-col gap-6">
            <div className="flex flex-wrap items-end gap-3">
              <div className="flex flex-col gap-1">
                <label className="text-xs font-semibold">Name</label>
                <input
                  value={addForm.name}
                  onChange={(e) => setAddForm({ ...addForm, name: e.target.value })}
                  required
                  className="border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                />
              </div>
              <div className="flex flex-col gap-1">
                <label className="text-xs font-semibold">Stock (g)</label>
                <input
                  type="number"
                  min="1"
                  value={addForm.amount_grams}
                  onChange={(e) =>
                    setAddForm({ ...addForm, amount_grams: e.target.value })
                  }
                  required
                  className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                />
              </div>
              <div className="flex flex-col gap-1">
                <label className="text-xs font-semibold">Spool Price (€)</label>
                <input
                  type="number"
                  min="0.01"
                  step="0.01"
                  value={addForm.total_price}
                  onChange={(e) =>
                    setAddForm({ ...addForm, total_price: e.target.value })
                  }
                  required
                  className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                />
              </div>
              <div className="flex flex-col gap-1">
                <label className="text-xs font-semibold">Cost/g</label>
                <span className="border-2 border-black bg-gray-100 px-3 py-1.5 text-sm">
                  €{costPerGram(addForm)}
                </span>
              </div>
              <div className="flex flex-col gap-1">
                <label className="text-xs font-semibold">Color</label>
                <input
                  type="color"
                  value={addForm.color_hex}
                  onChange={(e) =>
                    setAddForm({ ...addForm, color_hex: e.target.value })
                  }
                  required
                  className="h-[38px] w-16 border-2 border-black bg-white p-0"
                />
              </div>
              <div className="flex flex-col gap-1">
                <label className="text-xs font-semibold">Filament Type</label>
                <input
                  value={addForm.filament_type}
                  onChange={(e) =>
                    setAddForm({ ...addForm, filament_type: e.target.value })
                  }
                  required
                  placeholder="PLA"
                  className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                />
              </div>
            </div>

            <fieldset className="flex flex-col gap-3 border-2 border-black p-4">
              <legend className="px-2 text-xs font-semibold">
                Slicer profile (optional, defaults shown)
              </legend>
              <div className="flex flex-wrap gap-3">
                <div className="flex flex-col gap-1">
                  <label className="text-xs font-semibold">Nozzle Temp (°C)</label>
                  <input
                    type="number"
                    value={addForm.temperature}
                    onChange={(e) =>
                      setAddForm({ ...addForm, temperature: e.target.value })
                    }
                    placeholder="220"
                    className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                  />
                </div>
                <div className="flex flex-col gap-1">
                  <label className="text-xs font-semibold">Bed Temp (°C)</label>
                  <input
                    type="number"
                    value={addForm.bed_temperature}
                    onChange={(e) =>
                      setAddForm({ ...addForm, bed_temperature: e.target.value })
                    }
                    placeholder="60"
                    className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                  />
                </div>
                <div className="flex flex-col gap-1">
                  <label className="text-xs font-semibold">
                    First Layer Nozzle (°C)
                  </label>
                  <input
                    type="number"
                    value={addForm.first_layer_temperature}
                    onChange={(e) =>
                      setAddForm({
                        ...addForm,
                        first_layer_temperature: e.target.value,
                      })
                    }
                    placeholder="225"
                    className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                  />
                </div>
                <div className="flex flex-col gap-1">
                  <label className="text-xs font-semibold">
                    First Layer Bed (°C)
                  </label>
                  <input
                    type="number"
                    value={addForm.first_layer_bed_temperature}
                    onChange={(e) =>
                      setAddForm({
                        ...addForm,
                        first_layer_bed_temperature: e.target.value,
                      })
                    }
                    placeholder="65"
                    className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                  />
                </div>
                <div className="flex flex-col gap-1">
                  <label className="text-xs font-semibold">Diameter (mm)</label>
                  <input
                    type="number"
                    step="0.01"
                    value={addForm.filament_diameter}
                    onChange={(e) =>
                      setAddForm({ ...addForm, filament_diameter: e.target.value })
                    }
                    placeholder="1.75"
                    className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                  />
                </div>
                <div className="flex flex-col gap-1">
                  <label className="text-xs font-semibold">Extrusion Mult.</label>
                  <input
                    type="number"
                    step="0.01"
                    value={addForm.extrusion_multiplier}
                    onChange={(e) =>
                      setAddForm({
                        ...addForm,
                        extrusion_multiplier: e.target.value,
                      })
                    }
                    placeholder="1"
                    className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                  />
                </div>
                <div className="flex flex-col gap-1">
                  <label className="text-xs font-semibold">Density (g/cm³)</label>
                  <input
                    type="number"
                    step="0.01"
                    value={addForm.filament_density}
                    onChange={(e) =>
                      setAddForm({ ...addForm, filament_density: e.target.value })
                    }
                    placeholder="1.24"
                    className="w-28 border-2 border-black bg-white px-3 py-1.5 focus:outline-none"
                  />
                </div>
              </div>
            </fieldset>

            <div className="flex gap-3">
              <button
                type="submit"
                disabled={create.isPending}
                className="hover-black border-2 border-black bg-white hover:cursor-pointer px-4 py-1.5 font-semibold disabled:opacity-50"
              >
                Add
              </button>
              <button
                type="button"
                onClick={() => setAddForm(null)}
                className="border-2 border-black bg-white hover:cursor-pointer px-4 py-1.5 font-semibold hover:bg-gray-100"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}
    </div>
  )
}
